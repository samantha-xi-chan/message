package internal

import "os"

var (
	MONGO_URL = "mongodb://localhost:27017"
	AMQP_URL  = "amqp://guest:guest@localhost:5672/"
)

func init() {
	if os.Getenv("MONGO_SERVER") != "" {
		MONGO_URL = os.Getenv("MONGO_SERVER")
	}

	if os.Getenv("RABBITMQ_SERVER") != "" {
		AMQP_URL = os.Getenv("RABBITMQ_SERVER")
	}
}
