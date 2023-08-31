package sftp

type Config struct {
	User     string
	Password string
	Private  string
}

type Client struct {
	c Config
}

func New(cfg Config) *Client {
	cli := Client{
		c: cfg,
	}
	return &cli
}

func (c *client) Do() {

}

func (c *client) ReadDir(path string) {

}
