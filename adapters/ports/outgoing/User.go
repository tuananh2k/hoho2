package outgoing

import "hoho-framework-v2/model"

/*
create by: Hoangnd
create at: 2023-01-01
des: Xử lý thông tin hiển thị cho người dùng
*/
func ResponseUsers(us *model.User) string {

	return "Mr" + us.Name
}
