package core

type Auth struct {
	Username string `json:"username"`
	Secret   Secret `json:"secret"`
}
