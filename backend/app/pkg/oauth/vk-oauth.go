package oauth

import (
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/vk"
	"net/http"
)

func (oap *OAuthProvider) UseVKAuth(router *httprouter.Router) {
	goth.UseProviders(vk.New(
		oap.config.OAuth.VK.Key,
		oap.config.OAuth.VK.Secret,
		oap.config.OAuth.VK.CallbackUrl,
		"read",
	))

	router.POST("/api/oauth/vk/authorize", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		gothic.BeginAuthHandler(w, r)
	})

	router.POST("/api/oauth/vk/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		authUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		oap.OAuth(authUser, w, r)
	})
}
