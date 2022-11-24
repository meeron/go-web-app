package users

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id int `json:"id"`
}
