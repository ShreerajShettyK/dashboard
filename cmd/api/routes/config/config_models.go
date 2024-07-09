package config

type Config struct {
	MongoDBURI string `json:"MONGO_URI"`
	DBName     string `json:"DB_NAME"`
}
