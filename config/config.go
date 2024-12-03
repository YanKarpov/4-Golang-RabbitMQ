package config

// Config хранит настройки для подключения к RabbitMQ
type Config struct {
	Login    string
	Password string
	Host     string
	Port     string
}

// NewConfig возвращает новый объект Config с настройками
func NewConfig(login, password, host, port string) *Config {
	return &Config{
		Login:    login,
		Password: password,
		Host:     host,
		Port:     port,
	}
}
