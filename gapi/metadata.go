package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent" // works for API gateway
	xForwardedForHeader        = "x-forwarded-for"        // works for API gateway
	userAgentHeader            = "user-agent"             // works for grp client e.g. evans
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (srv *Server) extractMetadata(ctx context.Context) *Metadata {
	mtData := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtData.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtData.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtData.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtData.ClientIP = p.Addr.String()
	}

	return mtData
}
