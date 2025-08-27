package internal

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	tokenv1 "github.com/sazajun1390/tokentestserv/pkg/gen/token/v1"

	"crypto/sha256"

	models "github.com/sazajun1390/tokentestserv/pkg/bun/migrations"
	"github.com/uptrace/bun"
)

type UserService struct {
	db *bun.DB
}

func NewUserService(db *bun.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, req *connect.Request[tokenv1.CreateUserTokenRequest]) (*connect.Response[tokenv1.CreateUserTokenResponse], error) {

	passhash := sha256.Sum256([]byte(req.Msg.Password))
	_, err := s.db.NewInsert().Model(&models.User{
		UserEmail: req.Msg.UserEmail,
		Password:  req.Msg.Password,
	}).Exec(ctx)

	if err != nil {
		return nil, connect.NewError(connect.Internal, err)
	}

	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
