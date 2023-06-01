package host

type hostResolver struct{}

func NewHostResolver() (FetchHost, error) {
	return &hostResolver{}, nil
}

func (d *hostResolver) GetHost(domain string) (host Host, ok bool) {
	// TODO:解析/etc/hosts文件

	// 先读取本机/etc/hosts文件如果解析找不到则从dns中查找

	return
}

func (d *hostResolver) Close() error {
	return nil
}
