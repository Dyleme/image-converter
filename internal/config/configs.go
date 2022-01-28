package config

import (
	"os"
	"time"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/Dyleme/image-coverter/internal/rabbitmq"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// CollectiveConfig is a struct which contains configs for all needed constructors.
type CollectiveConfig struct {
	DB            *repository.DBConfig
	RabbitMQ      *rabbitmq.Config
	JWT           *jwt.Config
	AWS           *aws.Config
	AwsBucketName string
	LogLevel      string // string names from logrus.
	Port          string
}

// InitConfig is a function which create and initialize CollectiveConfig.
func InitConfig() (*CollectiveConfig, error) {
	db := &repository.DBConfig{
		UserName: os.Getenv("DBUSERNAME"),
		Password: os.Getenv("DBPASSWORD"),
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("DBSSLMODE"),
	}

	rabbitConfig := &rabbitmq.Config{
		User:     os.Getenv("RBUSER"),
		Password: os.Getenv("RBPASSWORD"),
		Host:     os.Getenv("RBHOST"),
		Port:     os.Getenv("RBPORT"),
	}

	ttl, err := time.ParseDuration(os.Getenv("TOKENTTL"))
	if err != nil {
		return nil, err
	}

	jwtConfig := &jwt.Config{
		SignedKey: os.Getenv("SIGNEDKEY"),
		TTL:       ttl,
	}

	port := os.Getenv("PORT")

	logLvl := os.Getenv("LOGLEVEL")

	awsBucketName := os.Getenv("AWS_BUCKET_NAME")
	awsConfig := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}

	return &CollectiveConfig{
		DB:            db,
		RabbitMQ:      rabbitConfig,
		JWT:           jwtConfig,
		Port:          port,
		AWS:           awsConfig,
		AwsBucketName: awsBucketName,
		LogLevel:      logLvl,
	}, nil
}
