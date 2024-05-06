package config

type Config struct {
	Debug      bool `yaml:"debug"`
	Workspaces struct {
		Exclude []string `yaml:"exclude"`
	} `yaml:"workspaces"`
}
