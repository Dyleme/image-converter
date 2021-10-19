package image

type User struct {
	ID       int    `json:"-"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
