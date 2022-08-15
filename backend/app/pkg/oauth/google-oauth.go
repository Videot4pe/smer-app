package oauth

import (
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"net/http"
	"os"
)

func (oap *OAuthProvider) UseGoogleAuth(router *httprouter.Router) {
	goth.UseProviders(google.New(
		os.Getenv(oap.config.OAuth.Google.Key),
		os.Getenv(oap.config.OAuth.Google.Secret),
		oap.config.OAuth.Google.CallbackUrl,
		"read",
	))

	router.POST("/api/oauth/google/authorize", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		gothic.BeginAuthHandler(w, r)
	})

	router.POST("/api/oauth/google/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		authUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		oap.OAuth(authUser, w, r)
	})
}
