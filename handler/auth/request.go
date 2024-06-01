package auth

type reqGoogleLogin struct {
	Ref string `query:"ref"`
}

type reqGoogleCallback struct {
	State string   `query:"state"`
	Code  string   `query:"code"`
	Scope []string `query:"scope"`
}

type reqDevSignUp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type reqDevSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
