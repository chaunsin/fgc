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
	c := client{
		c: cfg,
	}
	return &c
}

func (c *client) Do() {

}

func (c *client) ReadDir(path string) {

}
