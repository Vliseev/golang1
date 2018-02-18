package bench

type User struct {
	Browsers []string `json:"browsers,[]string"`
	Email    string `json:"email,string"`
	Name string `json:"name,string"`
}
