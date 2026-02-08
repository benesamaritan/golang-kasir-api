package handlers

import (
	"fmt"
	"net/http"
)

func WelcomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Selamat Datang Di API Aplikasi Kasir!")
	}
}
