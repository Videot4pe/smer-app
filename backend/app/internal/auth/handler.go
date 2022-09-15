package auth

import (
	"backend/internal/config"
	"backend/internal/domain/user"
	"backend/pkg/auth"
	"backend/pkg/logging"
	"backend/pkg/mailer"
	"backend/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Handler struct {
	logger  *logging.Logger
	storage *user.Storage
	ctx     context.Context
	cfg     *config.Config
}

type AuthenticatePayload struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type ChangePasswordPayload struct {
	Password string `json:"password"`
}

const (
	signinURL         = "/api/auth/signin"
	signupURL         = "/api/auth/signup"
	refreshURL        = "/api/auth/refresh"
	activateURL       = "/api/auth/activate/:hash"
	passwordResetURL  = "/api/auth/password-reset"
	changePasswordURL = "/api/auth/change-password"
)

func NewAuthHandler(ctx context.Context, storage *user.Storage, logger *logging.Logger, cfg *config.Config) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
		ctx:     ctx,
		cfg:     cfg,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.POST(signinURL, h.Signin)
	router.POST(signupURL, h.Signup)
	router.POST(refreshURL, h.Refresh)
	router.GET(activateURL, h.Activate)
	router.POST(passwordResetURL, h.PasswordReset)
	router.POST(changePasswordURL, h.ChangePassword)
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var credentials Credentials

	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, isVerified, err := h.storage.GetByCredentials(credentials.Email, credentials.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if !isVerified {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Not activated")
		return
	}

	jwtClaims := auth.AuthJwt{
		Data: auth.AuthJwtData{
			Id:    userId,
			Email: credentials.Email,
		},
	}

	token, err := auth.Encode(&jwtClaims, 10)

	// TODO email -> token ???
	refreshJwtClaims := auth.AuthJwt{
		Data: auth.AuthJwtData{
			Id:    userId,
			Email: token,
		},
	}

	refreshToken, err := auth.Encode(&refreshJwtClaims, 10)

	err = h.storage.UpdateRefreshToken(refreshToken, userId)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, AuthenticatePayload{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var payload AuthenticatePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error(err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.storage.IsRefreshTokenActual(payload.Token)
	if err != nil {
		h.logger.Error(err)
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	if userId == 0 {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Token is invalid")
		return
	}

	userInfo, err := h.storage.GetById(userId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwtClaims := auth.AuthJwt{
		Data: auth.AuthJwtData{
			Id:    userId,
			Email: userInfo.Email,
		},
	}
	payload.Token, err = auth.Encode(&jwtClaims, 10)

	refreshJwtClaims := auth.AuthJwt{
		Data: auth.AuthJwtData{
			Id:    userId,
			Email: payload.Token,
		},
	}
	refreshToken, err := auth.Encode(&refreshJwtClaims, 10)

	err = h.storage.UpdateRefreshToken(refreshToken, userId)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, AuthenticatePayload{
		Token:        payload.Token,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newUser user.User

	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &newUser); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO Transaction (create & send mail)
	userId, token, err := h.storage.Create(newUser, false)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	cfg := config.GetConfig()

	authMailerClient := GetMailerAuth(cfg, h.logger)
	activationLink := fmt.Sprintf("%v:%v/api/auth/activate/%v", cfg.Listen.ServerIP, cfg.Listen.Port, token)

	emailConfirmationParams := EmailConfirmationParams{
		Name:  newUser.Name,
		Email: newUser.Email,
		Link:  activationLink,
	}

	err = authMailerClient.SendMail(newUser.Email, "Email confirmation", EmailConfirmationTemplate, emailConfirmationParams)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Mail error")
		return
	}

	utils.WriteResponse(w, http.StatusOK, userId)
}

func (h *Handler) Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hash := ps.ByName("hash")
	err := h.storage.Activate(hash)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Activation error")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%v:%v/smers", h.cfg.Frontend.ServerIP, h.cfg.Frontend.Port), http.StatusTemporaryRedirect)
}

func (h *Handler) PasswordReset(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var email string

	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	email = string(body)

	userId, isVerified, err := h.storage.GetByEmail(email)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if !isVerified {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Not activated")
		return
	}

	token, err := h.storage.PasswordReset(userId)

	cfg := config.GetConfig()
	sender := mailer.SenderConfig{
		Host:     cfg.Mailer.Host,
		Port:     cfg.Mailer.Port,
		Username: cfg.Mailer.Username,
		Password: cfg.Mailer.Password,
	}
	mailClient := mailer.GetMailer(sender, h.logger)
	activationLink := fmt.Sprintf("%v:%v/change-password?token=%v", cfg.Frontend.ServerIP, cfg.Frontend.Port, token)

	mail := mailer.Mail{
		Username: email,
		Subject:  "Password reset",
		//Text:     result,
		Text: fmt.Sprintf("Pwd reset link: %v", activationLink),
	}
	err = mailClient.Send(mail)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Mail error")
		return
	}

	utils.WriteResponse(w, http.StatusOK, userId)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	token := queryValues.Get("token")
	var payload ChangePasswordPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Password change error: "+err.Error())
		return
	}

	err = h.storage.ChangePassword(token, payload.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Password change error: "+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
