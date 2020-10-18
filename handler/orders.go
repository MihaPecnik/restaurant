package handler

import (
	"encoding/json"
	"github.com/MihaPecnik/restaurant/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request models.Order
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.DB.CreateOrder(request)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, "Order created")
}

func (h *Handler) AddItemToOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderid, err := strconv.ParseInt(params["orderid"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	itemid, err := strconv.ParseInt(params["itemid"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.DB.AddItemToOrder(orderid, itemid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, "Item added to an order")
}

func (h *Handler) PayTheOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderid, err := strconv.ParseInt(params["orderid"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var funds float64
	err = json.NewDecoder(r.Body).Decode(&funds)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.DB.PayTheOrder(orderid, funds)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}
