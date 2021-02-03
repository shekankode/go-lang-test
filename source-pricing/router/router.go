package router

import (
	"source-pricing/controllers"
	"github.com/gorilla/mux"
	)

func Router() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/api/source-pricing/{product_id}/{store_id}", controllers.GetProductDetailHandler).Methods("GET")
    return router
}