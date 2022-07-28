package qdb

import (
	"errors"
	"github.com/spf13/viper"
)

var DefaultProfile string
var Profiles []ProfileDef

const defaultConfigFilename = "config"
const envPrefix = "qdb"
const qdbdirname = ".qdbctl"

type ProfileDef struct {
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
	DefaultProfile = viper.GetString("default-profile")
	if err := viper.UnmarshalKey("profiles", &Profiles); err != nil {
		return err
	}
	return nil
}

func SaveDefaultProfile(defaultProfileName string) error {
	viper.Set("default-profile", defaultProfileName)
	DefaultProfile = defaultProfileName
	return viper.WriteConfig()
}

func ProfileExists(name string) bool {
	for _, p := range Profiles {
		if p.Name == name {
			return true
		}
	}
	return false
}

func IsDefaultProfile(name string) bool {
	return DefaultProfile == name
}

func AddProfile(name string, url string) error {
	if ProfileExists(name) {
		return errors.New("profile '" + name + "' already exists")
	}
	profile := ProfileDef{Name: name, Url: url}
	Profiles = append(Profiles, profile)
	viper.Set("profiles", Profiles)
	if len(Profiles) == 1 {
		DefaultProfile = name
	}
	return nil
}

func DeleteProfile(name string) error {
	newProfiles := make([]ProfileDef, len(Profiles))
	found := false
	for _, p := range Profiles {
		if p.Name == name {
			found = true
		} else {
			newProfiles = append(newProfiles, p)
		}
	}
	if !found {
		return errors.New("profile '" + name + "' does not exist")
	}

	Profiles = newProfiles
	viper.Set("profiles", Profiles)
	if name == DefaultProfile {
		DefaultProfile = ""
	}
	return nil
}

func SaveConfig() error {
	return viper.WriteConfig()
}
