package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"rxdrag.com/entify/consts"
)

var c config

type config struct {
	v *viper.Viper
}

type DbConfig struct {
	Installed bool   `json:"installed"`
	Driver    string `json:"driver"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	Database  string `json:"database"`
	ID        int    `json:"id"`
}

const (
	TRUE  = "true"
	FALSE = "false"
)

const (
	PATH        = "."
	CONFIG_TYPE = "yaml"
	CONFIG_NAME = "config"
)

func init() {
	c.v = viper.New()
	c.v.SetConfigName(CONFIG_NAME) // name of config file (without extension)
	c.v.SetConfigType(CONFIG_TYPE) // REQUIRED if the config file does not have the extension in the name
	c.v.AddConfigPath(PATH)

	err := c.v.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		WriteConfig()
	}
}

func GetString(key string) string {
	return c.v.GetString(key)
}

func GetBool(key string) bool {
	return c.v.GetBool(key)
}

func GetInt(key string) int {
	return c.v.GetInt(key)
}

func SetString(key string, value string) {
	c.v.Set(key, value)
}

func SetInt(key string, value int) {
	c.v.Set(key, value)
}

func SetBool(key string, value bool) {
	c.v.Set(key, value)
}

func WriteConfig() {
	filePath := filepath.Join(PATH, CONFIG_NAME+"."+CONFIG_TYPE)
	err := c.v.WriteConfigAs(filePath)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Println("Can find config file and create:" + filePath)
}

func SetDbConfig(cfg DbConfig) {
	SetBool(consts.INSTALLED, cfg.Installed)
	SetString(consts.DB_DRIVER, cfg.Driver)
	SetString(consts.DB_DATABASE, cfg.Database)
	SetString(consts.DB_HOST, cfg.Host)
	SetString(consts.DB_PORT, cfg.Port)
	SetString(consts.DB_USER, cfg.User)
	SetString(consts.DB_PASSWORD, cfg.Password)
	SetInt(consts.ID, cfg.ID)
}

func GetDbConfig() DbConfig {
	var cfg DbConfig
	cfg.Installed = GetBool(consts.INSTALLED)
	cfg.Driver = GetString(consts.DB_DRIVER)
	cfg.Database = GetString(consts.DB_DATABASE)
	cfg.Host = GetString(consts.DB_HOST)
	cfg.Port = GetString(consts.DB_PORT)
	cfg.User = GetString(consts.DB_USER)
	cfg.Password = GetString(consts.DB_PASSWORD)
	cfg.ID = GetInt(consts.ID)
	return cfg
}

func ServiceId() int {
	if GetBool(consts.INSTALLED) {
		return GetInt(consts.ID)
	}
	return 0
}