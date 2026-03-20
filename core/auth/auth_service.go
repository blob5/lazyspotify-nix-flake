package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type AuthService struct {
	sptAuth *spotifyauth.Authenticator
	tknChannel chan *oauth2.Token
	authConfig *AuthConfig
}

func NewAuthService(authServer *AuthServer) *AuthService {
	authConfig := NewAuthConfig()
  sptAuth := spotifyauth.New(
		spotifyauth.WithRedirectURL(authServer.GetOauthRedirectURI()),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
      spotifyauth.ScopePlaylistReadPrivate,
      spotifyauth.ScopePlaylistModifyPublic,
      spotifyauth.ScopePlaylistModifyPrivate,
      spotifyauth.ScopePlaylistReadCollaborative,
      spotifyauth.ScopeUserFollowModify,
      spotifyauth.ScopeUserFollowRead,
      spotifyauth.ScopeUserLibraryModify,
      spotifyauth.ScopeUserLibraryRead,
      spotifyauth.ScopeUserReadCurrentlyPlaying,
      spotifyauth.ScopeUserReadPlaybackState,
      spotifyauth.ScopeUserModifyPlaybackState,
      spotifyauth.ScopeUserReadRecentlyPlayed,
      spotifyauth.ScopeUserTopRead,
      spotifyauth.ScopeUserReadEmail,
			),
		spotifyauth.WithClientID("565c1a413de9452da373f1ed3aa6afbe"),
		)
  return &AuthService{
    sptAuth: sptAuth,
    tknChannel: make(chan *oauth2.Token),
    authConfig: authConfig,
  }
}

func (a *AuthService) getAuthURL() string {
  return a.sptAuth.AuthURL(a.authConfig.state,
    oauth2.SetAuthURLParam("code_challenge_method", "S256"),
    oauth2.SetAuthURLParam("code_challenge", a.authConfig.codeChallenge),
    oauth2.SetAuthURLParam("client_id", "565c1a413de9452da373f1ed3aa6afbe"),
  )
}

func (a *AuthService) Authenticate(authServer *AuthServer) (*oauth2.Token, error) {
	authServer.InitAuthServer(a.makeOauthCallbackHandler())
	serverErrch := authServer.Start(a.authConfig)
	url := a.getAuthURL()
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	select {
  case tkn := <-a.tknChannel:
		authServer.Shutdown()
    return tkn, nil
  case err := <-serverErrch:
    return nil, err
	}
}

func (a *AuthService) makeOauthCallbackHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := a.sptAuth.Token(r.Context(), a.authConfig.state, r, oauth2.SetAuthURLParam("code_verifier", a.authConfig.codeVerifier))
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != a.authConfig.state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, a.authConfig.state)
		}
		a.tknChannel <- tok
	}
}

func (a *AuthService) GetSpotifyClient(tkn *oauth2.Token)  *spotify.Client {
	httpClient := a.sptAuth.Client(context.Background(), tkn)
	return spotify.New(httpClient)
}
