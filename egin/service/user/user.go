package user

type Info struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type User interface {
	Info(id int) (Info, error)
	Role(id int) ([]int, error)
}
