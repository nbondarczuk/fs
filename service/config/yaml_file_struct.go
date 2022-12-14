package config

type Document struct {
	Loggers    []EnvVarMapper `yaml:"loggers"`
	Servers    []EnvVarMapper `yaml:"servers"`
	SQLOptions []EnvVarMapper `yaml:"sqloptions"`
	Backends   []EnvVarMapper `yaml:"backends"`
	Stores     []EnvVarMapper `yaml:"stores"`
}

type EnvVarMapper struct {
	Kind string            `yaml:"kind"`
	Env  map[string]string `yaml:"env"`
}
