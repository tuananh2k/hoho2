package incoming

/*
create by: Hoangnd
create at: 2023-01-01
des: Khai báo và xử lý thông tin nhận từ request
*/
import (
	"errors"
	"hoho-framework-v2/model"
)

type UserParams struct {
	ID   string `json:"id" form:"id" param:"id" primary:"true" `
	Name string `json:"name" form:"name" param:"name"`
	Age  string `json:"age" form:"age" param:"age"`
}

func (uP UserParams) ValidateInput() error {
	if uP.ID == "" {
		return errors.New("Id can't be empty")
	}
	return nil
}
func (uP *UserParams) GetModel() model.User {
	return model.User{
		ID:   uP.ID,
		Name: uP.Name,
		Age:  uP.Age,
	}
}
