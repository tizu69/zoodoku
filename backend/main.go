package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"g.tizu.dev/zoodoku/backend/db"
	"g.tizu.dev/zoodoku/backend/gen/api/apiconnect"
	"github.com/rs/cors"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/validate"
)

//go:generate sqlc generate
//go:generate buf generate
//go:generate buf format -w

func main() {
	ctx := context.Background()

	dbPath := os.Getenv("ZOODOKU_DB")
	if dbPath == "" {
		dbPath = "zoodoku.db"
	}
	sqlDB, err := db.Open(ctx, dbPath)
	if err != nil {
		slog.Error("open database", "err", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	server := connect.NewServer(
		validate.NewServerInterceptor(),
	)

	apiconnect.RegisterGreetServiceHandler(server, &GreetServer{})
	apiconnect.RegisterUserServiceHandler(server, &UserServer{queries: db.New(sqlDB)})

	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://zoodoku.tizu.dev"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	p := new(http.Protocols)
	p.SetHTTP1(true)
	s := http.Server{
		Addr:      "localhost:8080",
		Handler:   c.Handler(mux),
		Protocols: p,
	}

	slog.Info("server started!", "addr", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		slog.Error("server failed", "err", err)
	}
}

func newAuthGate() connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			info, ok := connect.CallInfoForServerContext(ctx)
			if !ok || info.RequestHeader().Get("Authorization") == "" {
				return connect.NewError(connect.CodeUnauthenticated, "Unauthenticated")
			}
			return next(ctx, spec, stream)
		}
	}
}
