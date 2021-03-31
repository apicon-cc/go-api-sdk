package user

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
	Key      string `json:"-"`
	IP       string `json:"ip"`
}
