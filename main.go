package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

var db *sql.DB
var err error

func main() {

	//proof of concept only, do not store credentials in Github for production use
	//use token else and localized cred read file with encryption at rest
	db, err = sql.Open("mysql", "msandbox:msandbox@tcp(127.0.0.1:8033)/heb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//mux handler functions defined in functions.go
	router := mux.NewRouter()

	//reporting endpoints
	router.HandleFunc("/cart_items/{cust_id}", itemizedlist).Methods("GET")
	router.HandleFunc("/cart_items_tax/{cust_id}", itemizedtaxtotal).Methods("PUT")
	router.HandleFunc("/cart_items_full/{cust_id}", fulltaxtotal).Methods("PUT")
	router.HandleFunc("/cart_items_coupon/{cust_id}", couponAdd).Methods("PUT")

	//modifying cart endpoints
	router.HandleFunc("/cart_items_add/{cust_id}", itemAdd).Methods("POST")
	router.HandleFunc("/cart_items_remove/{cust_id}", itemRemove).Methods("POST")

	//assumes secure private network, not public plain text
	http.ListenAndServe(":8001", router)

}
