package models

type PlayRequest struct {
	Uri string `json:"uri"`
  SkipToUri string `json:"skip_to_uri"`
  Paused bool `json:"paused"`
}
