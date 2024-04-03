package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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
	EventSubsripter   *broker.EventSubsripter
}

func NewRouter(host string, port string, serviceuser service.UserServiceInterface, servicecontroller service.ControllerServiceInterface, broker *broker.EventSubsripter, hasher hasher.Interactor) *Router {
	r := &Router{UserService: serviceuser, ControllerService: servicecontroller, Hasher: hasher, EventSubsripter: broker}
	mux := http.NewServeMux()
	mux.HandleFunc("/login", r.login)
	mux.HandleFunc("/regsiter", r.register)
	mux.HandleFunc("/getidcontroller", r.getidcontroller)
	r.Server = &http.Server{
		Addr:    host + ":" + port,
		Handler: mux,
	}
	return r
}

func (s *Router) login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// проверка на валидность емайла
	var email = r.Header.Get("email")
	var password = r.Header.Get("password")
	user, err := s.UserService.GetByEmail(ctx, email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !s.Hasher.CompareAndHash(user.HashPassword, password) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = conn.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	user.HashPassword = ""
	data, err := json.Marshal(&user)
	conn.Write(ctx, websocket.MessageText, data)
	ct := conn.CloseRead(context.Background())
	var closed bool = true
	go func() {
		<-ct.Done()
		closed = false
		return

	}()
	ctrl, err := s.UserService.GetControllersByUserId(ctx, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var buffer chan []byte = make(chan []byte, 100)
	for i := 0; i < len(ctrl); i++ {
		go s.EventSubsripter.SubscribeMessange(ctx, strconv.Itoa(ctrl[i].Id_contorller), buffer)
	}
	go func() {
		for tmp := range buffer {
			var msg entitys.MessangeTypeZiroJson
			json.Unmarshal(tmp, &msg)
			var ans entitys.MessageFromFrontendJSON
			ans.Id = 801
			ans.Msgs = append(ans.Msgs, msg)
			tmp, _ := json.Marshal(&ans)
			conn.Write(ctx, websocket.MessageText, tmp)
		}
	}()
	for closed {
		msgtype, binry, err := conn.Read(ctx)
		if err != nil || msgtype != websocket.MessageText {
			conn.Write(ctx, websocket.MessageText, []byte("error sending data (type or connection)"))
			continue
		}
		var msg entitys.MessageFromFrontendJSON

		err = json.Unmarshal(binry, &msg)
		if err != nil {
			conn.Write(ctx, websocket.MessageText, []byte("ivalide json data"))
			continue
		}
		switch msg.Id {
		case 600:
			ansdata, err := s.ControllerService.GetListMessangesFromIdForUserId(ctx, msg.Rng.Count, msg.Rng.From, int(user.ID))
			if err != nil {
				conn.Write(ctx, websocket.MessageText, []byte("ivalide database select"))
				break
			}
			var ans entitys.MessageFromFrontendJSON
			ans.Id = 801
			ans.Msgs = ansdata
			bin, err := json.Marshal(&ans)
			if err != nil {
				conn.Write(ctx, websocket.MessageText, []byte("ivalide json parsing selected data from database"))
				break
			}
			conn.Write(ctx, websocket.MessageText, bin)
			break
		default:
			conn.Write(ctx, websocket.MessageText, []byte("error code json msg"))
			break
		}
	}
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
