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

	// need user login
	authLogin := r.PathPrefix("").Subrouter()
	authLogin.Use(middleware.AuthLoginMiddleware(app.Store))
	{
		userMethod := authLogin.PathPrefix("/user").Subrouter()
		{
			userMethod.HandleFunc("/current", app.UserHandler.GetCurrentUser).Methods("GET")
			userMethod.HandleFunc("/logout", app.UserHandler.UserLogout).Methods("POST")
			userMethod.HandleFunc("/update", app.UserHandler.UserUpdate).Methods("POST")
			userMethod.HandleFunc("/drop", app.UserHandler.UserDrop).Methods("POST")
			userMethod.HandleFunc("/search/tags", app.UserHandler.SearchUserByTags).Methods("GET")
			userMethod.HandleFunc("/recommend", app.UserHandler.RecommendUser).Methods("GET")
			userMethod.HandleFunc("/match", app.UserHandler.MatchUser).Methods("GET")
		}

		teamMethod := authLogin.PathPrefix("/team").Subrouter()
		{
			teamMethod.HandleFunc("/add", app.TeamHandler.AddTeam).Methods("POST")
			teamMethod.HandleFunc("/update", app.TeamHandler.UpdateTeam).Methods("POST")
			teamMethod.HandleFunc("/get", app.TeamHandler.GetTeamById).Methods("GET")
			teamMethod.HandleFunc("/list", app.TeamHandler.ListTeams).Methods("GET")
			teamMethod.HandleFunc("/join", app.TeamHandler.JoinTeam).Methods("POST")
			teamMethod.HandleFunc("/quit", app.TeamHandler.QuitTeam).Methods("POST")
			teamMethod.HandleFunc("/delete", app.TeamHandler.DeleteTeam).Methods("POST")
			teamMethod.HandleFunc("/list/my/create", app.TeamHandler.ListMyCreateTeams).Methods("GET")
			teamMethod.HandleFunc("/list/my/join", app.TeamHandler.ListMyJoinTeams).Methods("GET")
		}

	}

	return r
}
