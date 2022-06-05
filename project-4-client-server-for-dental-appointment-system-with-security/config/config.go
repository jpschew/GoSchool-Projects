// Package config includes a structure Config that stores user information and a function to read in all the values from a configuration file and return them as a structure.
package config

import (
	"github.com/spf13/viper"
)

// Config stores all information for user/admin for initialization purpose.
// the values are read by viper from a configuration file or environment variables
// viper uses mapstructure for unmarshalling, so we need to specify the variable name in environment file to map to the fields of structure
type Config struct {
	UserName  string `mapstructure:"UN"`
	FirstName string `mapstructure:"FN"`
	LastName  string `mapstructure:"LN"`
	Password  string `mapstructure:"PW"`
}

// LoadConfig reads in environment variables from a path provided if it exists.
// It will load in all the values from the environment file and return them as a struct.
func LoadConfig(path string, filename string) (Config, error) {

	var config Config

	// specify the location of config file to viper
	viper.AddConfigPath(path)

	// tell viper to look for a config file with specific name
	viper.SetConfigName(filename)

	// specify the type/format of file for viper to look for
	// in this case, it is env format
	viper.SetConfigType("env")

	// read values from env vars to viper
	viper.AutomaticEnv()

	// start reading in the config/env vars values
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// unmarshal decodes json into go struct
	// arg is the address of go struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}
	// fmt.Println(config, "config")

	return config, nil

}
