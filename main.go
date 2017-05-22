package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/mujz/ipp/config"
)

func NewServer() http.Handler {
	server := http.NewServeMux()
	server.HandleFunc("/", NotFoundHandler)
	server.HandleFunc("/current", AuthDecorator(CurrentHandler))
	server.HandleFunc("/next", AuthDecorator(NextHandler))

	server.HandleFunc("/login", LoginValidationDecorator(LoginHandler))
	server.HandleFunc("/signup", LoginValidationDecorator(SignupHandler))

	server.HandleFunc("/login/facebook", fbAuth.LoginHandler)
	server.HandleFunc("/login/facebook/callback", fbAuth.LoginCallbackHandler)
	return GlobalHeadersHandler(server)
}

func main() {
	port := config.Port
	server := NewServer()
	err := http.ListenAndServe(":"+port, handlers.CombinedLoggingHandler(os.Stdout, server))
	log.Fatal(err)
}
