package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// func writeToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Custom middleware")
// 		next.ServeHTTP(w, r)
// 	})
// }

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func sessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
