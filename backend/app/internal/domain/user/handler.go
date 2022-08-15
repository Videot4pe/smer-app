package user

import (
	"backend/pkg/auth"
	"backend/pkg/client/postgresql/model"
	"backend/pkg/logging"
	"backend/pkg/uploader"
	"backend/pkg/utils"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Handler struct {
	logger       *logging.Logger
	storage      *Storage
	filesStorage FilesStorage
	ctx          context.Context
}

type FilesStorage interface {
	Create(path string, name string) (uint16, error)
	GetById(id uint16) (string, error)
}

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:userId"
)

func NewUserHandler(ctx context.Context, storage *Storage, logger *logging.Logger, filesStorage FilesStorage) *Handler {
	return &Handler{
		filesStorage: filesStorage,
		logger:       logger,
		storage:      storage,
		ctx:          ctx,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.GET(usersURL, auth.RequireAuth(h.GetUser))
	router.PATCH(usersURL, auth.RequireAuth(h.UpdateUser))
	router.DELETE(usersURL, auth.RequireAuth(h.DeleteUser))
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pagination, err := model.NewPagination(r)

	users, err := h.storage.All(nil, pagination)
	if err != nil {
		h.logger.Error(err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, users)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := r.Context().Value("userId").(uint16)
	user, err := h.storage.GetById(id)

	if user.AvatarId != nil {
		avatarPath, err := h.filesStorage.GetById(*user.AvatarId)
		if err != nil {
			h.logger.Error(err)
			return
		}
		user.Avatar = &avatarPath
	}

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user User

	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, _, err := h.storage.Create(user, false)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusCreated, userId)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := r.Context().Value("userId").(uint16)

	var user User
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Avatar != nil && *user.Avatar != "" && user.AvatarId == nil {
		fileUploader := uploader.GetUploader(h.logger)
		path, name, err := fileUploader.Base64Upload(*user.Avatar)
		avatarId, err := h.filesStorage.Create(path, name)
		if err != nil {
			h.logger.Error(err)
			return
		}
		user.AvatarId = &avatarId
	}

	err = h.storage.Update(uint16(id), user)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, id)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("userId"), 16, 16)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.storage.Delete(uint16(id))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, id)
}
