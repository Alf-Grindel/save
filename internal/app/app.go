package app

import (
	"context"
	"github.com/Alf_Grindel/save/internal/handler"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"github.com/gorilla/sessions"
)

type Application struct {
	Logger      hlog.FullLogger
	Store       sessions.Store
	UserHandler *handler.UserHandler
}

func NewApplication() *Application {

	logger := hlog.DefaultLogger()
	ctx := context.Background()
	store := sessions.NewCookieStore([]byte("saveM814"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}

	userHandler := handler.NewUserHandler(ctx, store)

	app := &Application{
		Logger:      logger,
		Store:       store,
		UserHandler: userHandler,
	}

	return app
}
