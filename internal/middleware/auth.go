package middleware

import (
	"context"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/service"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"github.com/gorilla/sessions"
	"net/http"
)

func AuthLoginMiddleware(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, constant.UserLoginState)
			if err != nil {
				utils.RespWithErr(w, errno.SystemErr)
				return
			}
			auth, ok := session.Values["login"].(bool)
			if !ok || !auth {
				utils.RespWithErr(w, errno.NotLoginErr)
				return
			}
			account, ok := session.Values["user_account"].(string)
			if !ok || account == "" {
				utils.RespWithErr(w, errno.NotLoginErr)
				return
			}
			current, err := db.QueryUserByAccount(account)
			if err != nil {
				utils.RespWithErr(w, err)
				return
			}
			user := service.GetSafeUser(current)

			ctx := context.WithValue(r.Context(), constant.SessionKey, session)
			ctx = context.WithValue(ctx, constant.CtxUserInfoKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
