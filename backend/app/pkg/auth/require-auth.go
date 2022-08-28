package auth

import (
	"backend/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func RequireAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		bearer := strings.Replace(authHeader, "Bearer ", "", 1)
		_, claims, err := Decode(bearer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "userId", claims.Id))
		next(w, r, ps)
	}
}
