package model

// Database for database connection string
type Database struct {
	DataSource string `json:"dataSource" envconfig:"DB_DATASOURCE" required:"true"`
	Debug      bool   `json:"debug" envconfig:"DB_DEBUG"`
	Schema     string `json:"schema" envconfig:"DB_SCHEMA"`
	Type       string `json:"type" envconfig:"DB_TYPE" default:"sqlite3"`
}

// User to store the logged-in user details
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Token to check if user is logged-in
type Token struct {
	Role        string `json:"role"`
	Username    string `json:"username"`
	TokenString string `json:"token"`
}
