package config

import (
	"errors"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents all configurable values
type Config struct {
	// Either true for in-memory or false for MongoDB (default)
	InMemoryStorage bool `envconfig:"IN_MEMORY_STORAGE" required:"true" default:"false"`

	// The key used to sign outgoing and verify incoming JWTs
	JWTSigningKey string `envconfig:"JWT_SIGNING_KEY" required:"true" default:"super super secret secret"`

	// The version and build-date the binary.
	Version   string    `envconfig:"VERSION"`
	BuildDate time.Time `envconfig:"BUILD_DATE"`

	// What port to serve the api on (default 8080)
	APIPort int `envconfig:"API_PORT" required:"true" default:"8080"`

	// Rest-Layer debug mode
	Debug bool `envconfig:"DEBUG" required:"true" default:"true"`

	// Connection details to MongoDB
	MongoHost string `envconfig:"MONGO_HOST" default:"discapi-mongodb.svc.cluster.local"`
	MongoPort int    `envconfig:"MONGO_PORT" default:"27017"`
	MongoDB   string `envconfig:"MONGO_DB" default:"discapi"`
}

// ParseConfig takes env vars and returns golang values
func ParseConfig() (*Config, error) {
	ret := new(Config)
	if err := envconfig.Process("discapi", ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// Log outputs the configuration to the log
func (c *Config) Log() error {
	if c.JWTSigningKey == "" {
		return errors.New("JWT Signing Key cannot be empty string")
	}

	log.Println("==== Config ====")

	if c.InMemoryStorage {
		log.Println("Storage: In-memory")
	} else {
		log.Println("Storage: MongoDB")
		log.Printf("Mongo Host: %s\n", c.MongoHost)
		log.Printf("Mongo Port: %d\n", c.MongoPort)
		log.Printf("Mongo DB: %s\n", c.MongoDB)
	}

	log.Println("JWT Signing Key: --redacted--")

	if c.Debug {
		log.Println("Debug enabled")
	}

	log.Printf("Version: %s\n", c.Version)
	log.Printf("Build Date: %s\n", c.BuildDate.String())
	log.Printf("API Port: %d\n", c.APIPort)

	log.Println("==========")

	return nil
}
