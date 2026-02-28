package util

import "github.com/spf13/viper"

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// 告诉 viper 去哪个目录找配置文件
	viper.AddConfigPath(path)
	// 将很多配置信息事先写在 app.env 中
	viper.SetConfigName("app")
	viper.SetConfigType("env") // 也可以是其他文件类型，json，toml之类的

	// 启用环境变量读取（可覆盖同名配置）
	viper.AutomaticEnv()

	// 真正读取配置文件
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	// 把读取到的值解码到 Config 结构体返回
	err = viper.Unmarshal(&config)
	return
}