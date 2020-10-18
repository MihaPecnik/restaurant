package handler

import (
	"encoding/json"
	"github.com/MihaPecnik/restaurant/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var request models.ProductPrice
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.DB.CreateProduct(request)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, "Product created")
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.DB.ListProducts()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateThePrice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderid, err := strconv.ParseInt(params["itemid"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	var request models.ProductPrice
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.DB.UpdateThePrice(orderid, request.Cost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, "Price updated")
}