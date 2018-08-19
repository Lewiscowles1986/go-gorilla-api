package data

import (
    "github.com/satori/go.uuid"
)

func ParseUUID(str string) uuid.UUID {
    out, _ := uuid.FromString(str)
  	return out
}
