package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Homeppv2/api-go/internal/broker"
	"github.com/Homeppv2/api-go/internal/service"
	"github.com/Homeppv2/api-go/pkg/hasher"
	"github.com/Homeppv2/entitys"
	"nhooyr.io/websocket"
)

type Router struct {
	Server            *http.Server
	UserService       service.UserServiceInterface
	ControllerService service.ControllerServiceInterface
	Hasher            hasher.Interactor
}

var uriBroker = fmt.Sprintf("%s://%s:%s@%s:%s",
	os.Getenv("BROKER_PROTOCOL"),
	os.Getenv("BROKER_USERNAME"),
	os.Getenv("BROKER_PASSWORD"),
	os.Getenv("BROKER_HOST"),
	os.Getenv("BROKER_PORT"),
)

func NewRouter(host string, port string, serviceuser service.UserServiceInterface, servicecontroller service.ControllerServiceInterface, hasher hasher.Interactor) *Router {
	r := &Router{UserService: serviceuser, ControllerService: servicecontroller, Hasher: hasher}
	mux := http.NewServeMux()
	mux.HandleFunc("/login", r.login)
	mux.HandleFunc("/login/", r.login)

	mux.HandleFunc("/regsiter", r.register)
	mux.HandleFunc("/regsiter/", r.register)

	mux.HandleFunc("/getidcontroller", r.getidcontroller)
	r.Server = &http.Server{
		Addr:    host + ":" + port,
		Handler: mux,
	}
	return r
}

func (s *Router) login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var data []byte
	var login entitys.RequestAuth
	for {
		_, data, err = conn.Read(ctx)
		if err == nil {
			break
		}
	}
	err = json.Unmarshal(data, &login)
	if err != nil {
		conn.Close(http.StatusBadRequest, "user invalide "+err.Error())
		return
	}
	// проверка на валидность емайла
	var email = login.Email
	var password = login.Password
	log.Println(email)
	log.Println(password)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	user, err := s.UserService.GetByEmail(ctx, email)
	if err != nil {
		conn.Close(http.StatusBadRequest, "user invalide "+err.Error())
		return
	}
	if !s.Hasher.CompareAndHash(user.HashPassword, password) {
		conn.Close(http.StatusBadRequest, "user invalide")
		return
	}
	hp := user.HashPassword
	user.HashPassword = ""
	data, _ = json.Marshal(&user)
	err = conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return
	}
	ctrl, err := s.UserService.GetControllersByUserId(ctx, user.ID)
	if err != nil {
		conn.Close(http.StatusInternalServerError, "getting controller error "+err.Error())
		return
	}
	log.Println("количество контролеров у пользователя ", len(ctrl))
	var buffer chan []byte = make(chan []byte, 100)
	var gets []chan bool
	var eventSubsripter *broker.EventSubsripter
	eventSubsripter, err = broker.NewEventSubsripter(uriBroker, "current")
	if err != nil {
		log.Println("ошибка открытия uri " + uriBroker + " - " + err.Error())
		conn.Close(http.StatusInternalServerError, "rabbit error "+err.Error())
		return
	}
	for i := 0; i < len(ctrl); i++ {
		end := make(chan bool)
		err := eventSubsripter.SubscribeMessange(ctx, strconv.Itoa(ctrl[i].Id_contorller), buffer, end)
		if err != nil {
			log.Println("ошибка подписки на " + strconv.Itoa(ctrl[i].Id_contorller) + " - " + err.Error())
			continue
		}
		gets = append(gets, end)
	}
	go func() {
		for tmp := range buffer {
			log.Println("надо передать данные на пользователя ")
			var msg entitys.MessangeTypeZiroJson
			json.Unmarshal(tmp, &msg)
			log.Println(msg)
			/*
				if msg.RequestAuth == nil || !s.Hasher.CompareAndHash(hp, msg.RequestAuth.Password) || msg.RequestAuth.Email != user.Email {
					continue
				}
			*/
			var ans entitys.MessageFromFrontendJSON
			ans.Id = 801
			ans.Msgs = append(ans.Msgs, msg)
			tmp, _ := json.Marshal(&ans)
			conn.Write(ctx, websocket.MessageText, tmp)
		}
	}()
	for {
		msgtype, binry, err := conn.Read(ctx)
		if err != nil {
			break
		}
		if msgtype != websocket.MessageText {
			continue
		}
		var msg entitys.MessageFromFrontendJSON
		err = json.Unmarshal(binry, &msg)
		if err != nil {
			err = conn.Write(ctx, websocket.MessageText, []byte("ivalide json data"))
			if err != nil {
				break
			}
		}
		switch msg.Id {
		case 600:
			ansdata, err := s.ControllerService.GetListMessangesFromIdForUserId(ctx, msg.Rng.Count, msg.Rng.From, int(user.ID))
			if err != nil {
				conn.Write(ctx, websocket.MessageText, []byte("ivalide database select"))
			}
			var ans entitys.MessageFromFrontendJSON
			ans.Id = 801
			ans.Msgs = ansdata
			bin, err := json.Marshal(&ans)
			if err != nil {
				conn.Write(ctx, websocket.MessageText, []byte("ivalide json parsing selected data from database"))
			}
			err = conn.Write(ctx, websocket.MessageText, bin)
			if err != nil {
				break
			}
		default:
			err = conn.Write(ctx, websocket.MessageText, []byte("error code json msg"))
			if err != nil {
				break
			}
		}
	}
	for i := 0; i < len(gets); i++ {
		gets[i] <- true
		close(gets[i])
	}
	time.Sleep(5 * time.Second)
	close(buffer)
}

func (s *Router) register(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user entitys.User
	user.Email = r.Header.Get("email")
	user.Username = r.Header.Get("username")
	h, e := s.Hasher.HashPassword(r.Header.Get("password"))
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.HashPassword = h
	_, e = s.UserService.Register(ctx, user)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Router) getidcontroller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t, err := strconv.Atoi(r.Header.Get("type"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	n, err := strconv.Atoi(r.Header.Get("number"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.ControllerService.GetIdControllerByTypeAndNumber(r.Context(), t, n)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("idcontroller", strconv.Itoa(id))
}

func (r *Router) ListenAndServe() {
	r.Server.ListenAndServe()
}
