package helper

import "github.com/google/uuid"

func GetUuid() string {
	return uuid.Must(uuid.NewRandom()).String()
}
