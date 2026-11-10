package handler

import (
	"io"
	"net/http"
	"ragExample/internal/rag"
	"ragExample/internal/repo"

	"encoding/json"
)

type MyHandler struct {
	ServiceRAG *rag.Service
	repo       *repo.Repository
}

func MewMyHandler(servRAG *rag.Service, repo *repo.Repository) *MyHandler {
	return &MyHandler{
		ServiceRAG: servRAG,
		repo:       repo,
	}
}

func (h *MyHandler) StartQuery(w http.ResponseWriter, r *http.Request) {
	var data rag.QueryRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	w.Header().Set("Content-Type", "application/json")
	if err == io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(`{"error": "U must put question in request!"}`)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(`{"error": "Wrong parameters in body"}`)
		return
	}

	defer r.Body.Close()

	resQuery, err := h.ServiceRAG.Query(r.Context(), data.Question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resQuery)
}
