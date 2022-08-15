package auth

import (
	"backend/pkg/utils"
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func RequireAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		bearer := strings.Replace(authHeader, "Bearer ", "", 1)
		_, claims, err := DecodeJwt(bearer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "userId", claims.Id))
		next(w, r, ps)
	}
}
