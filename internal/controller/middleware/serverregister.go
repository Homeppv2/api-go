package middleware

import (
	"api-go/internal/controller"
	"api-go/internal/entity"
	"api-go/pkg/hasher"
	"context"
	"net/http"
)

type ServerRegister struct {
	Logf        func(f string, v ...interface{})
	Hasher      *hasher.Hasher
	UserService controller.UserService
}

func (s ServerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user entity.User
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
