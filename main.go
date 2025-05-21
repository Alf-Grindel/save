package main

import (
	"context"
	"github.com/Alf_Grindel/save/internal/dal"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"

	"github.com/Alf_Grindel/save/conf"
	. "github.com/Alf_Grindel/save/internal/app"
	"github.com/Alf_Grindel/save/internal/routes"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
)

func Init() {
	conf.LoadConfig()
	dal.Init()
}

func main() {
	Init()

	app := NewApplication()

	route := routes.SetUpRoutes(app)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE", "OPTION"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.ExposedHeaders([]string{"*"}),
		handlers.AllowCredentials(),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      cors(route),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	hlog.Info("start connect server on port 8080 ...")

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			hlog.Fatal("can not connect server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	hlog.Infof("Got signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
