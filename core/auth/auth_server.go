package auth

import (
	"fmt"
	"log"
	"net/http"
)

type AuthServer struct {
  host string
  port int
}

func NewAuthServer() *AuthServer {
  return &AuthServer{
    host: "127.0.0.1",
    port: 8287,
  }
}

func (authServer *AuthServer) GetOauthRedirectURI() string {
  return fmt.Sprintf("http://%s:%d/callback", authServer.host, authServer.port)
}

func (authServer *AuthServer) GetAuthServerAddress() string {
  return fmt.Sprintf("%s:%d", authServer.host, authServer.port)
}

// TODO: Check if port is occupied
// TODO: Check if callback func is registered
func (authServer *AuthServer) Start(authConfig *AuthConfig) {
  registerRoutes(authServer)
  startServer(authServer)
}

func (authServer *AuthServer) RegisterCallback(callbackFunc func(w http.ResponseWriter, r *http.Request)) {
  http.HandleFunc("/callback", callbackFunc)
}

func registerRoutes(authServer *AuthServer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
}

func startServer(authServer *AuthServer) {
	go http.ListenAndServe(authServer.GetAuthServerAddress(), nil)
}


