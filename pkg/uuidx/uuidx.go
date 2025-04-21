package uuidx

import (
	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"
	"github.com/google/uuid"
)

func Parse(str string) (uuid.UUID, error) {
	var id uuid.UUID
	var err error

	if id, err = uuid.Parse(str); err != nil {
		return uuid.Nil, apperrors.BadRequestHidden(err, "invalid uuid format")
	}

	return id, nil
}
