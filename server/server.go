package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	router *mux.Router
}

func (s *Server) Serve(port string) error {
	return http.ListenAndServe(port, s.router)
}

type Handler interface {
	CreateItem(w http.ResponseWriter, r *http.Request)
	ListItems(w http.ResponseWriter, r *http.Request)
	CreateOrder(w http.ResponseWriter, r *http.Request)
	AddItemToOrder(w http.ResponseWriter, r *http.Request)
	PayTheOrder(w http.ResponseWriter, r *http.Request)
	UpdateThePrice(w http.ResponseWriter, r *http.Request)
}


func NewServer(h Handler) *Server {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/item", h.CreateItem).Methods("POST")
	router.HandleFunc("/item", h.ListItems).Methods("GET")
	router.HandleFunc("/order", h.CreateOrder).Methods("POST")
	router.HandleFunc("/order/{orderid}/addItem/{itemid}", h.AddItemToOrder).Methods("POST")
	router.HandleFunc("/order/{orderid}/pay", h.PayTheOrder).Methods("PUT")
	router.HandleFunc("/item/{itemid}", h.UpdateThePrice).Methods("PUT")
	return &Server{
		router: router,
	}
}
