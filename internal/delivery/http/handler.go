package http

type Response struct {
	Status      any         `json:"status,omitempty"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}
