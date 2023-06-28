package redmine

// Redmine API config.
type Config struct {

	// Redmine host address.
	Host string `yaml:"host" mapstructure:"host"`

	// APIKey is a Redmine API key, it's obvious.
	APIKey string `yaml:"api_key" mapstructure:"api_key"`

	// Admin Login.
	Login string `yaml:"login" mapstructure:"login"`

	// Admin Password.
	Password string `yaml:"password" mapstructure:"password"`
}

// NewConfig returns empty config with default host.
func NewConfig() *Config {
	return &Config{}

}
