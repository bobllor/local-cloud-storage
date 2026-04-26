package api

// RequestUserLoginInfo is the request body containing the login information of the user.
type RequestUserLoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestUserRegisterInfo struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
