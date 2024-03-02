package server

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Config server config
type Config struct {
	Mode      string
	LogLevel  string       `yaml:"log_level"`
	ExpiredIn int          `yaml:"expired_in"` // redis 过期时间
	DB        *DBConfig    `yaml:"db"`
	Port      int          `yaml:"port"`
	Node      int64        `yaml:"node"`
	Redis     *RedisConfig `yaml:"redis"`
	Cos       *CosConfig   `yaml:"cos"`
}

// DBConfig config of db
type DBConfig struct {
	DataSources []string `yaml:"data_sources"`
	MaxCon      int      `yaml:"max_con"`
	MaxIdleCon  int      `yaml:"max_idle_con"`
	DriverName  string   `yaml:"driver_name"`
}

// RedisConfig redis config
type RedisConfig struct {
	Host string `yaml:"host"`
	Pwd  string `yaml:"pwd"`
}

// CosConfig cos config
type CosConfig struct {
	CosKey  string `yaml:"cos_key"`
	CosId   string `yaml:"cos_id"`
	CosAddr string `yaml:"cos_addr"`
}

func LoadLocalConfig(path, mode string) (*Config, error) {
	configPath := fmt.Sprintf("%s/%s/server.yaml", path, mode)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	data = []byte(os.ExpandEnv(string(data)))
	fmt.Println(string(data))

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
