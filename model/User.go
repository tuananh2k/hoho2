package model

type User struct {
	ID   string `json:"id" form:"id1" db:"id" primary:"true" `
	Name string `json:"name" form:"Name" db:"name"`
	Age  string `json:"age" form:"Age" db:"age"`
}

func (User) GetTableName() string { return "users" }
