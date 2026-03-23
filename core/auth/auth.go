package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	authServer  *AuthServer
	authService *AuthService
	keyring     *Keyring
}

const spotifyTokenKey = "token-v2"

func New() *Authenticator {
	authServer := NewAuthServer()
	return &Authenticator{
		authServer:  authServer,
		keyring:     NewSpotifyKeyring(),
		authService: NewAuthService(authServer),
	}
}

func (a *Authenticator) GetAuthToken(ctx context.Context) (*oauth2.Token, error) {
	tkn, err := a.keyring.GetToken(spotifyTokenKey)
	if err == nil {
		fmt.Printf("## token from keyring: %v\n", tkn.AccessToken)
		return tkn, nil
	}
	return a.ReAuthenticate(ctx)
}

func (a *Authenticator) GetClient(ctx context.Context) (*spotify.Client, error) {
	tkn, err := a.GetAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return a.authService.GetSpotifyClient(tkn), nil
}

func (a *Authenticator) ReAuthenticate(ctx context.Context) (*oauth2.Token, error) {
	fmt.Println("AuVthtenticating with spotify")
	tkn, err := a.authService.Authenticate(ctx, a.authServer)
	if err != nil {
		return nil, err
	}
	err = a.saveToken(tkn)
	if err != nil {
		log.Println("error saving token", err)
	}
	return tkn, nil
}

func (a *Authenticator) saveToken(token *oauth2.Token) error {
	return a.keyring.SetToken(spotifyTokenKey, token)
}
