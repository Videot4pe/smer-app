package smer

import (
	"backend/pkg/auth"
	"backend/pkg/client/postgresql/model"
	"backend/pkg/logging"
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
	logger  *logging.Logger
	storage *Storage
	ctx     context.Context
}

const (
	smersURL = "/api/smers"
	smerURL  = "/api/smers/:smerId"
)

func NewSmerHandler(ctx context.Context, storage *Storage, logger *logging.Logger) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
		ctx:     ctx,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.GET(smersURL, auth.RequireAuth(h.GetSmers))
	router.POST(smersURL, auth.RequireAuth(h.CreateSmer))
	router.GET(smerURL, auth.RequireAuth(h.GetSmer))
	router.PATCH(smerURL, auth.RequireAuth(h.UpdateSmer))
	router.DELETE(smerURL, auth.RequireAuth(h.DeleteSmer))
}

func (h *Handler) GetSmers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO -> записать в контекст сервиса

	userId := r.Context().Value("userId").(uint16)
	//h.storage.ctx = context.WithValue(h.storage.ctx, "userId", userId)

	// TODO withfilters
	// TODO withsorts
	pagination, err := model.NewPagination(r)
	sorts, err := model.NewSorts(r)
	filters, err := model.NewFilters(r)
	h.logger.Trace(filters)

	smers, meta, err := h.storage.All(userId, filters, pagination, sorts...)
	if err != nil {
		h.logger.Error(err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, utils.MetaData{
		Data: smers,
		Meta: meta,
	})
}

func (h *Handler) GetSmer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("smerId"), 16, 16)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	userId := r.Context().Value("userId").(uint16)
	smer, err := h.storage.GetById(userId, uint16(id))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, smer)
}

func (h *Handler) CreateSmer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var smer Smer
	userId := r.Context().Value("userId").(uint16)

	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &smer); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	smerId, err := h.storage.Create(smer, userId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusCreated, smerId)
}

func (h *Handler) UpdateSmer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("smerId"), 16, 16)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var smer Smer
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &smer); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := r.Context().Value("userId").(uint16)
	err = h.storage.Update(userId, uint16(id), smer)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, id)
}

func (h *Handler) DeleteSmer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("smerId"), 16, 16)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	userId := r.Context().Value("userId").(uint16)
	err = h.storage.Delete(userId, uint16(id))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, id)
}
