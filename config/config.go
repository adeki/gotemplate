package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ENV    string `envconfig:"ENV"`
	Server struct {
		Port string `envconfig:"SERVER_PORT"`
	}
	MySQL struct {
		User     string `envconfig:"MYSQL_USER"`
		Password string `envconfig:"MYSQL_PASSWORD"`
		Host     string `envconfig:"MYSQL_HOST"`
		Port     string `envconfig:"MYSQL_PORT"`
	}
	GCP struct {
		GoogleApplicationCredencials string `envconfig:"GOOGLE_APPLICATION_CREDENTIALS"`
	}
	AWS struct {
		Region string `envconfig:"AWS_REGION"`
		S3     struct {
			Buckets struct {
				PublicBucket  string `envconfig:"AWS_S3_BUCKETS_PUBLIC_BUCKET"`
				PrivateBucket string `envconfig:"AWS_S3_BUCKETS_PRIVATE_BUCKET"`
				MediaBucket   string `envconfig:"AWS_S3_BUCKETS_MEDIA_BUCKET"`
			}
			Endpoint       string `envconfig:"AWS_S3_ENDPOINT"`
			DisableSSL     bool   `envconfig:"AWS_S3_DISABLESSL"`
			ForcePathStyle bool   `envconfig:"AWS_S3_FORCEPATHSTYLE"`
		}
		CloudFront struct {
			DistributionID string `envconfig:"AWS_CLOUDFRONT_DISTRIBUTION_ID"`
		}
		MediaConvert struct {
			Role string `envconfig:"AWS_MEDIACONVERT_ROLE"`
		}
	}
}

var (
	cfg  Config
	once sync.Once
)

func Load() Config {
	once.Do(func() {
		cfg = process()
	})
	return cfg
}

func process() Config {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatalln(err)
	}
	return c
}

func (c Config) MySQLDatabase1Dsn() string {
	return c.mysqlDsn("database1")
}

func (c Config) MySQLDatabase2Dsn() string {
	return c.mysqlDsn("database2")
}

func (c Config) mysqlDsn(database string) string {
	return fmt.Sprint(c.MySQL.User, ":", c.MySQL.Password, "@(", c.MySQL.Host, ":", c.MySQL.Port, ")/", database, "?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo")
}
