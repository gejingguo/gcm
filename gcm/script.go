package gcm

// 脚本
type Script struct {
	ID int `json:"id"`			// 编号
	Name string	`json:"name"`	// 名称
	Desc string	`json:"desc"`	// 简介
	Body string `json:"body"`	// 内容
}
