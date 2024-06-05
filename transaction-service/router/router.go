package router

import (
    "github.com/gorilla/mux"
    "transaction-service/handlers"
    "net/http"
)

func SetupRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/transactions", handlers.CreateTransaction).Methods("POST")
    r.HandleFunc("/transactions/{id}", handlers.GetTransaction).Methods("GET")
    r.HandleFunc("/transactions/{id}/pay", handlers.PayTransaction).Methods("POST")

    // HTML pages
    r.HandleFunc("/create", handlers.RenderCreateTransactionPage).Methods("GET")
    r.HandleFunc("/pay/{id}", handlers.RenderPayTransactionPage).Methods("GET")
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    return r
}
