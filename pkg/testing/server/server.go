package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StopFn func()

func RunServer(t *testing.T, server abstractions.Server, port int) (*grpc.ClientConn, StopFn) {
	var err error

	go func() {
		err = server.Start()
		require.NoError(t, err)
	}()

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), dialOptions...)
	require.NoError(t, err)

	return conn, func() {
		err = conn.Close()
		require.NoError(t, err)

		stopCtx, cancel := context.WithTimeout(t.Context(), time.Second*5)
		defer cancel()

		err = server.Shutdown(stopCtx)
		require.NoError(t, err)
	}
}
