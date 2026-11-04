package main

import (
	"log/slog"
	"net/http"

	"g.tizu.dev/zoodoku/backend/gen/api/apiconnect"
	"github.com/rs/cors"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/validate"
)

func main() {
	server := connect.NewServer(
		validate.NewServerInterceptor(),
	)

	apiconnect.RegisterGreetServiceHandler(server, &GreetServer{})

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
