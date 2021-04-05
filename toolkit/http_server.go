package toolkit

import (
	//	"math/rand"
	"net/http"
	"time"

)

func MyListenAndServe(addr string, handler http.Handler) error {
	server := &http.Server{
		ReadTimeout: time.Second * 5,
		WriteTimeout: time.Second * 10,
		ReadHeaderTimeout: time.Second * 3,
		Addr: addr, 
		Handler: handler,
	}
	return server.ListenAndServe()
}
	
func MyListenAndServeTLS(addr string, handler http.Handler, cert, key string) error {
	server := &http.Server{
		ReadTimeout: time.Second * 5,
		WriteTimeout: time.Second * 10,
		ReadHeaderTimeout: time.Second * 3,
		Addr: addr, 
		Handler: handler,
	}
	return server.ListenAndServeTLS(cert, key)
}	
	