package tern

import (
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/stretchr/testify/require"
)

func RunMigration(t *testing.T, connString, migrationPath string) {
	conn, err := pgx.Connect(t.Context(), connString)
	require.NoError(t, err)

	migrator, err := migrate.NewMigrator(t.Context(), conn, "migrations")
	require.NoError(t, err)

	err = migrator.LoadMigrations(migrationPath)
	require.NoError(t, err)

	err = migrator.Migrate(t.Context())
	require.NoError(t, err)
}
