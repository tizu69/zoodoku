package main

import (
	"context"
	"fmt"

	"g.tizu.dev/zoodoku/backend/gen/api"
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *api.GreetRequest) (*api.GreetResponse, error) {
	return &api.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}
