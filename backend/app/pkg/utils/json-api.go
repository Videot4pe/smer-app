package utils

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int16  `json:"status"`
	Title  string `json:"title"`
}

type Meta struct {
	TotalItems uint64 `json:"totalItems"`
	TotalPages uint64 `json:"totalPages"`
}

type MetaData struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&JsonResponse{Data: data}); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: int16(errorCode), Title: errorMsg}})
}
