package mattermost

import "time"

// Mattermost API config.
type Config struct {
	UserEmail              string        `yaml:"user_email" mapstructure:"user_email"`
	UserPassword           string        `yaml:"user_password" mapstructure:"user_password"`
	UserName               string        `yaml:"user_name" mapstructure:"user_name"`
	TeamName               string        `yaml:"team_name" mapstructure:"team_name"`
	ServerAddress          string        `yaml:"server_address" mapstructure:"server_address"`
	ServerWsAddress        string        `yaml:"server_ws_address" mapstructure:"server_ws_address"`
	TimeBeforeNotification time.Duration `yaml:"time_before_notification" mapstructure:"time_before_notification"`
	ReportFileName         string        `yaml:"report_file_name" mapstructure:"report_file_name"`
	Messages
}
type Messages struct {
	Help         string `yaml:"help"`
	NoReport     string `yaml:"noreport"`
	Start        string `yaml:"start"`
	Finish       string `yaml:"finish"`
	ReportRemind string `yaml:"reportremind"`
	NoApiKey     string `yaml:"noapikey"`
}

// NewMattemost creates new MattemostConfig instance
// with no default configurations.
func NewConfig() (mt *Config) {
	mt = new(Config)
	return
}
func NewMessagesConfig() (ms *Messages) {
	ms = new(Messages)
	return
}
