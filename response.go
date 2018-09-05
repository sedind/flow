package app

// Response represents response object for API
type Response struct {
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
	Success bool        `json:"success"`
}
