package qdb

import (
	"errors"
	"github.com/spf13/viper"
)

var DefaultConnectionName string
var ConnectionDefs []ConnectionDef

const defaultConfigFilename = "config"
const envPrefix = "qdb"
const qdbdirname = ".qdbctl"

type ConnectionDef struct {
	Name string
	Url  string
}

func LoadConfig() error {
	viper.SetConfigName(defaultConfigFilename)

	viper.AddConfigPath("$HOME/.qdbctl")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigType("yaml")

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	viper.AutomaticEnv()
	DefaultConnectionName = viper.GetString("default-connection")
	if err := viper.UnmarshalKey("connections", &ConnectionDefs); err != nil {
		return err
	}
	return nil
}

func ConnectionByName(connectionName string) (ConnectionDef, error) {
	for _, c := range ConnectionDefs {
		if connectionName == c.Name {
			return c, nil
		}
	}
	return ConnectionDef{}, errors.New("connection '" + connectionName + "' does not exist")
}

func SetAsDefaultConnection(connectionName string) {
	DefaultConnectionName = connectionName
	viper.Set("default-connection", connectionName)
}

func ConnectionExists(name string) bool {
	for _, p := range ConnectionDefs {
		if p.Name == name {
			return true
		}
	}
	return false
}

func IsDefaultConnection(name string) bool {
	return DefaultConnectionName == name
}

func AddConnection(name string, url string) error {
	if ConnectionExists(name) {
		return errors.New("connection '" + name + "' already exists")
	}
	connection := ConnectionDef{Name: name, Url: url}
	ConnectionDefs = append(ConnectionDefs, connection)
	viper.Set("connections", ConnectionDefs)
	if len(ConnectionDefs) == 1 {
		SetAsDefaultConnection(name)
	}
	return nil
}

func DeleteConnection(connName string) error {
	newConnection := make([]ConnectionDef, len(ConnectionDefs))
	found := false
	for _, p := range ConnectionDefs {
		if p.Name == connName {
			found = true
		} else {
			newConnection = append(newConnection, p)
		}
	}
	if !found {
		return errors.New("connection '" + connName + "' does not exist")
	}

	ConnectionDefs = newConnection
	viper.Set("connections", ConnectionDefs)
	if connName == DefaultConnectionName {
		DefaultConnectionName = ""
	}
	return nil
}

func SaveConfig() error {
	return viper.WriteConfig()
}
