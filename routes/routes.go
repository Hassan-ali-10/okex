package routes

import (
	"net/http"
    "github.com/gorilla/mux"
    api "okex/routes/api"
    config "okex/config"

	
)

func SetRouters() {
	
	  r := mux.NewRouter()
	  r.HandleFunc("/executeOrders", api.CreateBook).Methods("POST")
      http.ListenAndServe(":"+config.PORT, r)
}