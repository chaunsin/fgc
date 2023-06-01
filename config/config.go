package config

type Config struct {
	Url string `json:"url"`
	Dir string `json:"dir"`
}

func New(dir string) *Config {
	return &Config{
		Url: "114.217.31.201:37101",
		Dir: dir,
	}
}


