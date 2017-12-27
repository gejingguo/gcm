package gcm


// 主机信息
type Host struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Host string	`json:"host"`		// format ip:port
	//Port int
	User string	`json:"user"`
	Pass string `json:"pass"`
	RootPath string `json:"root_path"`
}



