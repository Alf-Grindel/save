package app

import (
	"github.com/Alf_Grindel/save/conf"
	"github.com/Alf_Grindel/save/internal/handler"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
)

type Application struct {
	Logger      hlog.FullLogger
	Store       *redistore.RediStore
	UserHandler *handler.UserHandler
	TeamHandler *handler.TeamHandler
}

func NewApplication() *Application {

	logger := hlog.DefaultLogger()
	store, err := redistore.NewRediStore(constant.MaxIdleNum, constant.TCP, conf.Redis.Addr, "", "", []byte(constant.SessionKey))
	if err != nil {
		hlog.Fatal(err)
	}
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
	}

	userHandler := handler.NewUserHandler(store)
	teamHandler := handler.NewTeamHandler()

	app := &Application{
		Logger:      logger,
		Store:       store,
		UserHandler: userHandler,
		TeamHandler: teamHandler,
	}

	return app
}
