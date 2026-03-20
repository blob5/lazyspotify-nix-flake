package auth

import (
	"fmt"

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

func (a *Authenticator) Authenticate() (*oauth2.Token, error) {
  tkn, err := a.keyring.GetToken("token")
  if err == nil {
    return tkn, nil
  }
  fmt.Println("Authtenticating with spotify")
	tkn = a.authService.Authenticate(a.authServer)
	a.saveToken(tkn)
	return tkn, nil
}

func (a *Authenticator) GetClient() (*spotify.Client,error) {
  tkn, err := a.Authenticate()
  if err != nil {
    return nil, err
  }
  return a.authService.GetSpotifyClient(tkn), nil
}

func (a *Authenticator) saveToken(token *oauth2.Token) {
	err := a.keyring.SetToken("token", token)
	if err != nil {
		fmt.Println("Error saving token:", err)
	}
}
