package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)


type Keyring struct {
	service string
}

func NewSpotifyKeyring() *Keyring {
  return &Keyring{service: "spotify"}
}

func (k *Keyring) GetString(key string) (string, error) {
  return keyring.Get(k.service, key)
}

func (k *Keyring) SetString(key, value string) error {
	return keyring.Set(k.service, key, value)
}

func (k *Keyring) GetToken(key string) (*oauth2.Token, error) {
	savedToken := oauth2.Token{}
  ser_token, err := k.GetString(key)
	if(err != nil) {
    return nil, err
  }
  err = json.Unmarshal([]byte(ser_token), &savedToken)
	if (err != nil){
		fmt.Println(err)
		return nil , err
	}
	return &savedToken, nil
}

func (k *Keyring) SetToken(key string, token *oauth2.Token) error {
  ser_token, err := json.Marshal(token)
  if (err != nil){
    fmt.Println(err)
    return err
  }
  return k.SetString(key, string(ser_token))
}
