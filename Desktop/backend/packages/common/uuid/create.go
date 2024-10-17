package uuid

import guuid "github.com/google/uuid"

func Create() string {
	return guuid.Must(guuid.NewRandom()).String()
}
