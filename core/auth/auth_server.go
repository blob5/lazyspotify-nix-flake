package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthServer struct {
  host string
  port int
	httpServer *http.Server
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


func (authServer *AuthServer) Start(authConfig *AuthConfig) chan error {
  return startServer(authServer)
}


func (authServer *AuthServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Second)
  defer cancel()
  if authServer.httpServer != nil {
    err := authServer.httpServer.Shutdown(ctx)
    if err != nil {
      return err
    }
  }
  return nil
}

func (authServer *AuthServer) InitAuthServer(oauthRedirectCallbackFunc func(w http.ResponseWriter, r *http.Request)) {
	mux := http.NewServeMux()
	registerRoutes(mux,oauthRedirectCallbackFunc)
	server := &http.Server{
  	Addr: authServer.GetAuthServerAddress(),
	}
  authServer.httpServer = server
}

func registerRoutes(mux *http.ServeMux, oauthRedirectCallbackFunc func(w http.ResponseWriter, r *http.Request)) {
  mux.HandleFunc("/callback", oauthRedirectCallbackFunc)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
}


func startServer(authServer *AuthServer) chan error {
	errCh := make(chan error)
	go func() {
		err := authServer.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()
  return errCh
}


