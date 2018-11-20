package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
)

type Joker struct {
	Name  string `json:"name" gorm:"primary_key"`
	Jokes int    `json:"jokes"`
}

type Handler struct {
	db *gorm.DB
}

func main() {
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	var h Handler
	// db, err := gorm.Open("sqlite3", "jokers.db")
	db, err := gorm.Open("postgres", dbURL)
	defer db.Close()
	if err != nil {
		log.Fatal("unable to connect to DB", err)
	}
	h.db = db
	h.db.AutoMigrate(&Joker{})
	r := mux.NewRouter()
	// create the login route based on the api-attempt!
	r.HandleFunc("/add/{joker}", h.addJoker).Methods("GET")
	r.HandleFunc("/joke/{joker}", h.addToJoker).Methods("GET")
	r.HandleFunc("/score", h.score).Methods("GET")
	r.HandleFunc("/score", h.addMany).Methods("POST")
	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(fs)

	log.Fatal(http.ListenAndServe(port, r))
}

func (h Handler) addMany(w http.ResponseWriter, r *http.Request) {
	var jokers []Joker
	err := json.NewDecoder(r.Body).Decode(&jokers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, joker := range jokers {
		var jk Joker
		if err = h.db.Where("name=?", joker.Name).First(&jk).Error; gorm.IsRecordNotFoundError(err) {
			h.db.Create(&joker)
		} else {
			jk.Jokes = joker.Jokes
			h.db.Save(&jk)
		}
	}
}

func (h Handler) addToJoker(w http.ResponseWriter, r *http.Request) {
	var joker Joker
	vars := mux.Vars(r)
	name := vars["joker"]
	h.db.Where("name=?", name).First(&joker)
	joker.Jokes++
	h.db.Save(&joker)
	err := json.NewEncoder(w).Encode(joker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) addJoker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["joker"]
	joker := Joker{Name: name, Jokes: 0}
	h.db.Create(&joker)
	err := json.NewEncoder(w).Encode(joker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) score(w http.ResponseWriter, r *http.Request) {
	var jokers []Joker
	h.db.Find(&jokers)
	err := json.NewEncoder(w).Encode(jokers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
