package gapi

import (
	"context"
	"database/sql"

	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/blessedmadukoma/go-simple-bank/pb"
	"github.com/blessedmadukoma/go-simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := srv.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "passwords do not match")
	}

	accessToken, accessPayload, err := srv.tokenMaker.CreateToken(user.Username, srv.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := srv.tokenMaker.CreateToken(
		user.Username,
		srv.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	session, err := srv.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		// UserAgent:    ctx.Request.UserAgent(),
		// ClientIp:     ctx.ClientIP(),
		UserAgent: "",
		ClientIp:  "",
		IsBlocked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "failed to create session: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	return rsp, nil
}
