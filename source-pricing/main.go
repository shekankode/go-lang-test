package main

import (
	"net/http"
	"source-pricing/router"
)


func main() {
	
	router := router.Router()
	http.ListenAndServe(":8000", router)
	
}

