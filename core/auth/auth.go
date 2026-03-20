package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type Authenticator struct{
	authServer *AuthServer
	authService *AuthService
	keyring *Keyring
}

func New() *Authenticator {
  authServer := NewAuthServer()
  return &Authenticator{
    authServer: authServer,
    keyring: NewSpotifyKeyring(),
    authService: NewAuthService(authServer),
  }
}

func (a *Authenticator) Authenticate(ctx context.Context) (*oauth2.Token, error) {
  tkn, err := a.keyring.GetToken("token")
  if err == nil {
    return tkn, nil
  }
  fmt.Println("Authtenticating with spotify")
	tkn,err = a.authService.Authenticate(ctx,a.authServer)
	if err != nil {
		return nil, err
	}
	err = a.saveToken(tkn)
  if err != nil {
    log.Println("error saving token", err)
  }
	return tkn, nil
}

func (a *Authenticator) GetClient(ctx context.Context) (*spotify.Client,error) {
  tkn, err := a.Authenticate(ctx)
  if err != nil {
    return nil, err
  }
  return a.authService.GetSpotifyClient(tkn), nil
}

func (a *Authenticator) saveToken(token *oauth2.Token) error {
	return a.keyring.SetToken("token", token)
}
