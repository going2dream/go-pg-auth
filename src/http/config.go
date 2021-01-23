package http

type Config struct {
	Environment string `yaml:"environment"`
	BindIP      string `yaml:"bind_ip"`
	BindPort    string `yaml:"bind_port"`
}
