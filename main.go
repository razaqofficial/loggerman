package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"loggerman/controller/HomeController"
	_ "loggerman/database/model/userModel"
	"net/http"
)

const  (
	STATIC_DIR = "/static/"
	PORT       = "8080"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	router.HandleFunc("/", HomeController.Index).Methods("GET")
	router.HandleFunc("/about", HomeController.About).Methods("GET")
	router.HandleFunc("/loggers", HomeController.Loggers).Methods("GET")
	router.HandleFunc("/add/logger", HomeController.AddLogger).Methods("GET","POST")
	router.HandleFunc("/forum", HomeController.Forum).Methods("GET")
	router.HandleFunc("/contact", HomeController.Contact).Methods("GET","POST")


	log.Fatal(http.ListenAndServe(":"+PORT,router))
}


func openConnection() *gorm.DB {
	db,err := gorm.Open("mysql","")
	if err != nil {
		panic(err.Error())
	}
	return db
}
