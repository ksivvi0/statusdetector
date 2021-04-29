package config

//Config - структура настроек проекта
type Config struct {
	ListenAddress    string //IPv4 адрес для прослушивания
	ListenPort       uint16 //порт для прослушивания
	LogPath          string //путь для файла журнала
	ConnectionString string //строка подключения с БД
}

func NewConfig(address string, port uint16, logPath, connString string) (*Config, error) {
	return &Config{
		ListenAddress:    address,
		ListenPort:       port,
		LogPath:          logPath,
		ConnectionString: connString,
	}, nil
}
