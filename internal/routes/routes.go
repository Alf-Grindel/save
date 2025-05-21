package routes

import (
	"github.com/Alf_Grindel/save/internal/app"
	"github.com/Alf_Grindel/save/internal/middleware"
	"github.com/gorilla/mux"
)

func SetUpRoutes(app *app.Application) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.JsonContentTypeMiddleware)

	public := r.PathPrefix("").Subrouter()
	{
		public.HandleFunc("/user/register", app.UserHandler.UserRegister).Methods("POST")
		public.HandleFunc("/user/login", app.UserHandler.UserLogin).Methods("POST")
	}

	authLogin := r.PathPrefix("").Subrouter()
	authLogin.Use(middleware.AuthLoginMiddleware(app.Store))
	{
		authLogin.HandleFunc("/user/current", app.UserHandler.GetCurrentUser).Methods("GET")
		authLogin.HandleFunc("/user/logout", app.UserHandler.UserLogout).Methods("POST")
		authLogin.HandleFunc("/user/update", app.UserHandler.UserUpdate).Methods("POST")
		authLogin.HandleFunc("/user/drop", app.UserHandler.UserDrop).Methods("POST")
		authLogin.HandleFunc("/user/search/tags", app.UserHandler.SearchUserByTags).Methods("GET")
	}

	return r
}
