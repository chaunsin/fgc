package config

type Config struct {
	Url string `json:"url"`
	Dir string `json:"dir"`
}

func New(dir string) *Config {
	return &Config{
		Url: "0.0.0.0:8080",
		Dir: dir,
	}
}
