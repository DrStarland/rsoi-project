package middleware

import (
	"net/http"

	"gateway/pkg/myjson"
	"gateway/pkg/session"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func Auth(next httprouter.Handle, sm session.SessionsManager, logger *zap.SugaredLogger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logger.Infoln("TOKETOKETOKE ", r.Header.Get("Authorization"))

		sess, err := sm.Check(r)
		if err != nil {
			if err == session.ErrNoAuth {
				myjson.JsonError(w, http.StatusUnauthorized, err.Error())
				return
			}
			myjson.JsonError(w, http.StatusUnauthorized, err.Error())
			return
		} else if sess == nil {
			myjson.JsonError(w, http.StatusUnauthorized, "")
			return
		}

		r.Header.Set("X-User-Name", sess.User.Username)
		ctx := session.ContextWithSession(r.Context(), sess)
		next(w, r.WithContext(ctx), ps)
	}
}
