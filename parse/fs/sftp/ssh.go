package sftp

type Config struct {
	User     string
	Password string
	Private  string
}

type client struct {
	c Config
}

func New(cfg Config) *client {
	cli := client{
		c: cfg,
	}
	return &cli
}

func (c *client) Do() {

}

func (c *client) ReadDir(path string) {

}
