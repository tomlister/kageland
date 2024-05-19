package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomlister/kageland/handlers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.HandleFunc("/api/top", handlers.TopHandler)
	router.HandleFunc("/api/new", handlers.NewHandler)
	router.HandleFunc("/api/search", handlers.SearchHandler).Methods("GET")
	router.HandleFunc("/api/shader", handlers.ShaderPostHandler).Methods("POST")
	router.HandleFunc("/api/shader", handlers.ShaderGetHandler).Methods("GET")
	router.HandleFunc("/api/like", handlers.ShaderLikeHandler).Methods("POST")
	router.HandleFunc("/api/unlike", handlers.ShaderUnlikeHandler).Methods("POST")
	router.HandleFunc("/api", handlers.IndexHandler)
	router.ServeHTTP(w, r)
}
