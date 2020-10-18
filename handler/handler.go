package handler

import (
	"encoding/json"
	"github.com/MihaPecnik/restaurant/models"
	"net/http"
)

type Database interface {
	CreateProduct(request models.ProductPrice) error
	ListProducts() ([]models.ProductPrice, error)
	CreateOrder(order models.Order) error
	AddItemToOrder(orderId, itemId int64) error
	PayTheOrder(orderId int64, payment float64) (models.Receipt,error)
	UpdateThePrice(id int64, cost float64) error
}

type Handler struct {
	DB Database
}

func NewHandler(database Database) *Handler {
	return &Handler{
		DB: database,
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
