package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/blessedmadukoma/go-simple-bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

// authorizeUser checks if the request metadata contains a valid token
func (srv *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization token")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)

	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	if strings.ToLower(fields[0]) != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := fields[1]
	payload, err := srv.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	return payload, nil
}
