package image

type User struct {
	Id       int    `json:"-"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
