package config

type Options struct {
	Server struct {
		Host string
		Port int
	}
	Logger struct {
		Level      string
		Encoding   string
		OutputPath string
	}
	App struct {
		Environment string
		BasePath    string
	}
}
