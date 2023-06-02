package builder

type Options struct {
	Mode        string // 读取文件方式 local:默认 sftp ftp
	OrgName     string // 组织名称
	OrderName   string // 排序节点名称
	ChannelName string // 排序组织名称
	User        string // 用户名
	Pem         bool   // 证书生成的格式 false:pem文件格式(默认) true:路径方式
	DoubleTls   bool   // 生成tls false:单tls(默认) true:双tls
	CA          bool   // 是否开启ca false:关闭(默认) true:开启
	Metrics     bool   // 是否生成Metrics false:关闭(默认) true:开启
	Operations  bool   // 是否生成Operations false:关闭(默认) true:开启
	FileType    string // 生成文件类型 yaml(默认) json

	Language string
}
