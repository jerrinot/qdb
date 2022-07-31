package config

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
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"user"`
	Password string `json:"password"`
	IsCloud  bool   `json:"cloud"`
}

func LoadConfig() error {
	viper.SetConfigName(defaultConfigFilename)
	viper.AddConfigPath("$HOME/.qdbctl")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigType("yaml")
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

func AddConnection(conn ConnectionDef) error {
	if ConnectionExists(conn.Name) {
		return errors.New("connection '" + conn.Name + "' already exists")
	}
	ConnectionDefs = append(ConnectionDefs, conn)
	viper.Set("connections", ConnectionDefs)
	if len(ConnectionDefs) == 1 {
		SetAsDefaultConnection(conn.Name)
	}
	return nil
}

func DeleteConnection(connName string) error {
	newConnection := make([]ConnectionDef, 0)
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
