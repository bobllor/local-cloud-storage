package api

// RequestUserInfo is the request body containing the login information of the user.
type RequestUserInfo struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
