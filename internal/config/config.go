package config

type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	BotToken string `yaml:"botToken" env-required:"true"`
	ChatID string `yaml:"chatID" env-required:"true"`
}