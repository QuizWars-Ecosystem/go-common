package uuidx

import (
	"github.com/gofrs/uuid"
	pgxuuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

func NewUUIDFromString(id string) (pgxuuid.UUID, error) {
	res, err := uuid.FromString(id)
	if err != nil {
		return pgxuuid.UUID{}, err
	}
	return pgxuuid.UUID{UUID: res}, nil
}
