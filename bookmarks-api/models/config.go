package models

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Mode string `json:"mode"`
	DB   struct {
		BoltDB string `yaml:"boltdb"`
	} `yaml:"db"`
}

func (conf *Config) IsProduction() bool {
	return conf.Mode == "production"
}
