package config

var Local = Config{
	Environment: Environment{Type: "LOCAL"},
	Database: Database{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "user",
		Pass:     "password",
		Database: "goluca",
		SSLMode:  "disable",
	},
	HTTPServer: HTTPServer{
		Host: "localhost",
		Port: "3333",
	},
}
