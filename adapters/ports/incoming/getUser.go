package incoming

/*
create by: Hoangnd
create at: 2023-01-01
des: Khai báo và xử lý thông tin nhận từ request
*/

type GetUserApi struct {
	Page     string `json:"id" form:"id" param:"id" primary:"true" `
	PageSize string `json:"name" form:"name" param:"name"`
	Seearch  string `json:"age" form:"age" param:"age"`
	Ids      string `json:"age" form:"age" param:"age"`
}
