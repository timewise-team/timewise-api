package storage

import (
	"os"
)

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}
