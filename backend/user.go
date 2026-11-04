package main

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"g.tizu.dev/zoodoku/backend/db"
	"g.tizu.dev/zoodoku/backend/gen/api"
	"g.tizu.dev/zoodoku/backend/gen/api/apiconnect"
)

const (
	// friendcodeAlphabet deliberately leaves out 0/O and 1/I/L
	friendcodeAlphabet    = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	friendcodeLength      = 8
	maxFriendcodeAttempts = 5
)

type UserServer struct {
	queries *db.Queries
}

var _ apiconnect.UserServiceHandler = (*UserServer)(nil)

func (s *UserServer) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name must not be blank"))
	}

	for range maxFriendcodeAttempts {
		friendcode, err := newFriendcode()
		if err != nil {
			slog.ErrorContext(ctx, "generate friendcode", "err", err)
			return nil, connect.NewError(connect.CodeInternal, errors.New("could not create user"))
		}

		user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
			Name:       name,
			Friendcode: friendcode,
		})
		if err == nil {
			return &api.CreateUserResponse{
				Name:       user.Name,
				Friendcode: user.Friendcode,
			}, nil
		}
		if !isUniqueViolation(err) {
			slog.ErrorContext(ctx, "create user", "err", err)
			return nil, connect.NewError(connect.CodeInternal, errors.New("could not create user"))
		}
	}

	slog.ErrorContext(ctx, "ran out of friendcode attempts", "name", name)
	return nil, connect.NewError(connect.CodeInternal, errors.New("could not create user"))
}

func newFriendcode() (string, error) {
	var buf [friendcodeLength]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	for i, b := range buf {
		buf[i] = friendcodeAlphabet[int(b)%len(friendcodeAlphabet)]
	}
	return string(buf[:]), nil
}

func isUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
