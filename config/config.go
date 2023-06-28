package config

import (
	"strings"

	"github.com/spf13/viper"

	psql "mattermost/pkg/bun"
	"mattermost/pkg/logger"
	"mattermost/pkg/mattermost"
	"mattermost/pkg/redmine"
)

// Config of entire app.
type Config struct {
	Mattermost *mattermost.Config   `yaml:"mattermost" mapstructure:"mattermost"`
	DB         *psql.Config         `yaml:"db" mapstructure:"db"`
	Logger     *logger.Config       `yaml:"logger" mapstructure:"logger"`
	Redmine    *redmine.Config      `yaml:"redmine" mapstructure:"redmine"`
	Messages   *mattermost.Messages `yaml:"messageText" mapstructure:"messageText"`
}

// New creates a new config instance with
// no default configurations.
func New() (conf *Config) {
	conf = new(Config)
	conf.Mattermost = mattermost.NewConfig()
	conf.DB = psql.NewConfig()
	conf.Redmine = redmine.NewConfig()
	conf.Logger = logger.NewConfig()
	conf.Messages = mattermost.NewMessagesConfig()
	return
}

// Load the Config from configuration file.
// This method panics on error.
func (c *Config) Load() *Config {

	var v = viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	v.AutomaticEnv()

	var err error

	if err = v.Unmarshal(c); err != nil {
		panic(err)
	}

	return c
}
