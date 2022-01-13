package handlers

import (
	"net/http"
)

type IndexData struct {
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Kageland API"))
}
