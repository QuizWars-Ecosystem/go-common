package clients

import (
	"context"
	"net"
	"time"

	"github.com/exaring/otelpgx"
	"go.opentelemetry.io/otel/trace"

	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresClient(ctx context.Context, url string, options *PostgresOptions) (*pgxpool.Pool, error) {
	var opts *PostgresOptions

	if options == nil {
		opts = NewPostgresOptions(url)
	} else {
		opts = options
	}

	pingCtx, pingCtxCancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer pingCtxCancel()

	client, err := pgxpool.NewWithConfig(ctx, opts.Config)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	if opts.traceEnabled {
		_ = otelpgx.RecordStats(client)
	}

	if err = client.Ping(pingCtx); err != nil {
		return nil, apperrors.Internal(err)
	}

	return client, nil
}

type PostgresOptions struct {
	*pgxpool.Config
	traceEnabled bool
}

func NewPostgresOptions(url string) *PostgresOptions {
	opts := &PostgresOptions{}
	opts.Config, _ = pgxpool.ParseConfig(url)

	return opts
}

func (o *PostgresOptions) WithHost(host string) *PostgresOptions {
	o.ConnConfig.Host = host
	return o
}

func (o *PostgresOptions) WithUsername(username string) *PostgresOptions {
	o.ConnConfig.User = username
	return o
}

func (o *PostgresOptions) WithPassword(password string) *PostgresOptions {
	o.ConnConfig.Password = password
	return o
}

func (o *PostgresOptions) WithDatabase(database string) *PostgresOptions {
	o.ConnConfig.Database = database
	return o
}

func (o *PostgresOptions) WithDialFunc(fn func(ctx context.Context, network, addr string) (net.Conn, error)) *PostgresOptions {
	o.ConnConfig.DialFunc = fn
	return o
}

func (o *PostgresOptions) WithConnectTimeout(time time.Duration) *PostgresOptions {
	o.ConnConfig.ConnectTimeout = time
	return o
}

func (o *PostgresOptions) WithHealthCheckPeriod(time time.Duration) *PostgresOptions {
	o.HealthCheckPeriod = time
	return o
}

func (o *PostgresOptions) WithMinCons(amount int32) *PostgresOptions {
	o.MinConns = amount
	return o
}

func (o *PostgresOptions) WithMaxCons(amount int32) *PostgresOptions {
	o.MaxConns = amount
	return o
}

func (o *PostgresOptions) WithConnMaxLifetime(time time.Duration) *PostgresOptions {
	o.MaxConnLifetime = time
	return o
}

func (o *PostgresOptions) WithTracerProvider(provider trace.TracerProvider) *PostgresOptions {
	o.traceEnabled = true
	o.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithTracerProvider(provider),
	)
	return o
}
