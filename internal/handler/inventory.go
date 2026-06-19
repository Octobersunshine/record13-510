package handler

import (
	"encoding/json"
	"inventory-service/internal/service"
	"net/http"
	"strings"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *InventoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/stock", h.handleStock)
	mux.HandleFunc("/api/stock/sku/", h.handleSKUStock)
	mux.HandleFunc("/api/stock/set/", h.handleSetStock)
	mux.HandleFunc("/api/stock/all", h.handleAllStock)
}

func (h *InventoryHandler) handleStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if strings.HasPrefix(id, "set_") {
		h.getSetStock(w, id)
	} else {
		h.getSKUStock(w, id)
	}
}

func (h *InventoryHandler) handleSKUStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/stock/sku/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "sku id is required")
		return
	}

	h.getSKUStock(w, id)
}

func (h *InventoryHandler) handleSetStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/stock/set/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "set id is required")
		return
	}

	h.getSetStock(w, id)
}

func (h *InventoryHandler) handleAllStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stocks, err := h.svc.ListAllStock()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, stocks)
}

func (h *InventoryHandler) getSKUStock(w http.ResponseWriter, skuID string) {
	stock, err := h.svc.GetSKUStock(skuID)
	if err != nil {
		if err == service.ErrSKUNotFound {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeSuccess(w, stock)
}

func (h *InventoryHandler) getSetStock(w http.ResponseWriter, setID string) {
	stock, err := h.svc.GetSetStock(setID)
	if err != nil {
		if err == service.ErrSetNotFound {
			writeError(w, http.StatusNotFound, err.Error())
		} else if err == service.ErrSKUNotFound {
			writeError(w, http.StatusBadRequest, "set contains missing sku: "+err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeSuccess(w, stock)
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
	})
}
