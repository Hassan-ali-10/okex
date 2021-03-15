package routes

import (
	"net/http"
    "github.com/gorilla/mux"
    api "okex/routes/api"
    config "okex/config"
    

	
)

func SetRouters() {
	
	  r := mux.NewRouter()
	  r.HandleFunc("/A", api.Home).Methods("GET")
	  r.HandleFunc("/executeOrders", api.ExecuteOrdersPostRequest).Methods("GET")
      http.ListenAndServe(":"+config.PORT, r)
}