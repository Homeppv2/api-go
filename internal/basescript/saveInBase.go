package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Homeppv2/api-go/internal/broker"
	"github.com/Homeppv2/api-go/internal/database"
	"github.com/Homeppv2/api-go/internal/service"
	"github.com/Homeppv2/entitys"
	"github.com/jackc/pgx/v5/pgxpool"
)

var uriBroker = fmt.Sprintf("%s://%s:%s@%s:%s",
	os.Getenv("BROKER_PROTOCOL"),
	os.Getenv("BROKER_USERNAME"),
	os.Getenv("BROKER_PASSWORD"),
	os.Getenv("BROKER_HOST"),
	os.Getenv("BROKER_PORT"),
)

func main() {

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	l.Info("success initializing logger")

	// db
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("POSTGRESQL_USER"), os.Getenv("POSTGRESQL_PASSWORD"), os.Getenv("POSTGRESQL_HOST"), os.Getenv("POSTGRESQL_PORT"), os.Getenv("POSTGRESQL_BASE")))
	if err != nil {
		l.Info("failed to connect to postgresql: %s", err.Error())
		return
	}
	defer func(db *pgxpool.Pool) {
		db.Close()
	}(db)

	err = db.Ping(context.Background())
	if err != nil {
		l.Info("failed to ping to postgresql: %s", err.Error())
		return
	}
	base := database.NewDatabase(db)
	service := service.NewControllerService(base)
	l.Info("success connecting to postgresql")

	sb, err := broker.NewEventSubsripter(uriBroker, "base")
	if err != nil {
		log.Println("ошибка открытия uri " + uriBroker + " - " + err.Error())
		return
	}
	var ctrls []entitys.ControllersData
	var buffer chan []byte = make(chan []byte, 300)
	go func() {
		for tmp := range buffer {
			log.Println("надо сохранить данные в базу ")
			var msg entitys.MessangeTypeZiroJson
			json.Unmarshal(tmp, &msg)
			log.Println(msg)
			/*
				if msg.RequestAuth == nil {
					continue
				}
			*/
			if msg.One != nil {
				service.CreateMessageTypeOne(context.Background(), *msg.One)
			}
			if msg.Two != nil {
				service.CreateMessageTypeTwo(context.Background(), *msg.Two)
			}
			if msg.Three != nil {
				service.CreateMessageTypeThree(context.Background(), *msg.Three)
			}
		}
	}()
	var ends []chan bool
	ctrls, err = base.GetListContorllers(context.Background())
	log.Println(len(ctrls))
	for i := 0; i < len(ctrls); i++ {
		var end = make(chan bool)
		sb.SubscribeMessange(context.Background(), strconv.Itoa(ctrls[i].Id_contorller), buffer, end)
		ends = append(ends, end)
	}
	// t := time.NewTicker(3 * time.Minute)
	finish := make(chan bool)
	go func() {
		for {
			select {
			/*
				case <-t.C:
					var ends2 []chan bool
					ctrls, err = base.GetListContorllers(context.Background())
					for i := 0; i < len(ctrls); i++ {
						var end = make(chan bool)
						sb.SubscribeMessange(context.Background(), strconv.Itoa(ctrls[i].Id_contorller), buffer, end)
						ends2 = append(ends2, end)
					}
					for i := 0; i < len(ends); i++ {
						ends[i] <- true
						close(ends[i])
					}
					log.Println(len(ctrls))
					ends = ends2
			*/
			case <-finish:
				for i := 0; i < len(ends); i++ {
					ends[i] <- true
					close(ends[i])
				}
				time.Sleep(5 * time.Second)
				close(buffer)
				break
			}
		}
	}()
	l.Debug("success starting application")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit
	finish <- true

	l.Info("application has been shut down")

}
