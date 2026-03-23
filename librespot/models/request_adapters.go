package models

import "encoding/json"

func NewPlayRequest(uri string, skipToUri string, paused bool) ([]byte, error) { 
	playRequest:= PlayRequest{
    Uri: uri,
    SkipToUri: skipToUri,
    Paused: paused,
  }
	playRequestJson, err := json.Marshal(playRequest)
	if err != nil {
		return nil,err
	}
  return playRequestJson,nil
}
