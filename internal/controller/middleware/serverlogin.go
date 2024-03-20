package middleware

import (
	"api-go/internal/controller"
	"api-go/internal/entity"
	"api-go/internal/infrastructure"
	"api-go/pkg/hasher"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"nhooyr.io/websocket"
)

type ServerLogin struct {
	Logf             func(f string, v ...interface{})
	EventSubsripter  *infrastructure.EventSubsripter
	Hasher           *hasher.Hasher
	UserService      controller.UserService
	ControlerService controller.ControllerService
	// SocketGateway  *gateway.SocketGateway
}

func (s ServerLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	ctrl, err := s.ControlerService.GetControllersByUserId(ctx, user.ID)
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
			var msg entity.MessangeTypeZiroJson
			json.Unmarshal(tmp, &msg)
			var ans entity.MessageFromFrontendJSON
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
		var msg entity.MessageFromFrontendJSON

		err = json.Unmarshal(binry, &msg)
		if err != nil {
			conn.Write(ctx, websocket.MessageText, []byte("ivalide json data"))
			continue
		}
		switch msg.Id {
		case 600:
			ansdata, err := s.ControlerService.GetCountMessangesFromIdForUserId(ctx, msg.Rng.Count, msg.Rng.From, int(user.ID))
			if err != nil {
				conn.Write(ctx, websocket.MessageText, []byte("ivalide database select"))
				break
			}
			var ans entity.MessageFromFrontendJSON
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
