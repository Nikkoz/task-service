package config

type Environment string

const (
	EnvLocal      Environment = "local"
	EnvProduction Environment = "production"
)

func (env Environment) IsLocal() bool {
	return env == EnvLocal
}

func (env Environment) IsProduction() bool {
	return env == EnvProduction
}
