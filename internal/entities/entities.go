package entities

// EnvConfig represents the configuration structure for the application.
type EnvConfig struct {
	Debug            bool     `default:"true" split_words:"true"`                 // Flag indicating debug mode (default: true)
	Port             int      `default:"8080" split_words:"true"`                 // Port for server to listen on (default: 8080)
	Db               Database `split_words:"true"`                                // Database configuration
	AcceptedVersions []string `required:"true" split_words:"true" default:"v1.0"` // List of accepted API versions (required)
	JWTSecretKey     string   `split_words:"true"`
}

// Database represents the configuration for the database connection.
type Database struct {
	Driver    string `split_words:"true"`
	User      string // Database username
	Password  string // Database password
	Port      int    // Database port
	Host      string // Database host
	DATABASE  string // Database name
	Schema    string // Database schema
	MaxActive int    // Maximum number of active connections
	MaxIdle   int    // Maximum number of idle connections
}

type MetaData struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	Next        int `json:"next"`
	Prev        int `json:"prev"`
}
