package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const(
   Smongo = "mongo"
   Spostgres = "postgres"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `json:"logger"`
	DB     DBConf     `json:"db"`
	Server ServConf   `json:"server"`
}

type LoggerConf struct {
	Level string `json:"level"`
	LogFile string `json:"logFile"`
}

type DBConf struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
	Timeout int `json:"timeout"`
	NumOpenConns int `json:"numOpenConns"`
	ConnLifeTime int `json:"connLifeTime"`
}

type ServConf struct {
	Port int `json:"port"`
	ReadTimeOut int `json:"readTimeOut"`
	WriteTimeOut int `json:"writeTimeOut"`
}

func NewConfig(cfgPath string) Config {
	var data []byte
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("Error read config file: %s", err.Error()))
	}

    var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Error get configs: %s", err.Error()))
	}

	return cfg
}