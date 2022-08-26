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
	Hash     string `json:"hash"`
	Password string `json:"password"`
}

const (
	signinURL               = "/api/auth/signin"
	signupURL               = "/api/auth/signup"
	refreshURL              = "/api/auth/refresh"
	activateURL             = "/api/auth/activate/:hash"
	passwordResetURL        = "/api/auth/password-reset"
	passwordResetWebhookURL = "/api/auth/password-reset/:hash"
	changePasswordURL       = "/api/auth/change-password"
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
	router.GET(passwordResetWebhookURL, h.PasswordResetWebhook)
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

	jwtClaims := auth.NewJwtClaims(credentials.Email, userId)
	token, err := jwtClaims.EncodeJwt()

	refreshJwtClaims := auth.NewJwtClaims(token, userId)
	refreshToken, err := refreshJwtClaims.EncodeJwt()

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

	jwtClaims := auth.NewJwtClaims(userInfo.Email, userId)
	payload.Token, err = jwtClaims.EncodeJwt()

	refreshJwtClaims := auth.NewJwtClaims(payload.Token, userId)
	refreshToken, err := refreshJwtClaims.EncodeJwt()

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
	userId, hash, err := h.storage.Create(newUser, false)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	cfg := config.GetConfig()

	authMailerClient := GetMailerAuth(cfg, h.logger)
	activationLink := fmt.Sprintf("%v:%v/api/auth/activate/%v", cfg.Listen.ServerIP, cfg.Listen.Port, hash)

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
	http.Redirect(w, r, fmt.Sprintf("%v/smers", h.cfg.Listen.ServerIP), http.StatusTemporaryRedirect)
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

	// Transaction (create & send mail)
	userId, isVerified, hash, err := h.storage.GetByEmailAndGenerateHash(email)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if !isVerified {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Not activated")
		return
	}

	cfg := config.GetConfig()
	sender := mailer.SenderConfig{
		Host:     cfg.Mailer.Host,
		Port:     cfg.Mailer.Port,
		Username: cfg.Mailer.Username,
		Password: cfg.Mailer.Password,
	}
	mailClient := mailer.GetMailer(sender, h.logger)
	activationLink := fmt.Sprintf("%v:%v/change-password/%v", cfg.Frontend.ServerIP, cfg.Frontend.Port, hash)

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

func (h *Handler) PasswordResetWebhook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO Implement
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var payload ChangePasswordPayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	err = h.storage.PasswordReset(payload.Hash, payload.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Password reset error")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%v/signin", h.cfg.Listen.ServerIP), http.StatusTemporaryRedirect)
}
