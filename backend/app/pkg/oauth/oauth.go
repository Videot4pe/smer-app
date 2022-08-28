package oauth

import (
	"backend/internal/config"
	"backend/internal/domain/user"
	"backend/pkg/auth"
	"backend/pkg/logging"
	"backend/pkg/utils"
	"net/http"

	"github.com/markbates/goth"
)

type OAuthProvider struct {
	logger  *logging.Logger
	config  *config.Config
	storage *user.Storage
}

func GetOAuthProvider(logger *logging.Logger, cfg *config.Config, storage *user.Storage) *OAuthProvider {
	return &OAuthProvider{
		logger:  logger,
		config:  cfg,
		storage: storage,
	}
}

func (oap *OAuthProvider) OAuth(authUser goth.User, w http.ResponseWriter, r *http.Request) {

	newUser := user.User{
		Username: authUser.NickName,
		Name:     authUser.FirstName,
		Surname:  authUser.LastName,
		Email:    authUser.Email,
	}
	userId, _, err := oap.storage.Create(newUser, true)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwtClaims := &auth.AuthJwt{
		Email: newUser.Email,
		Id:    userId,
	}
	token, err := jwtClaims.Encode(10)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, token)
}
