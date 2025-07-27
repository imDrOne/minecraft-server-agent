package http_server

import (
	"net"
)

type HttpOption func(*HttpServer)

func WithPort(value string) HttpOption {
	return func(s *HttpServer) {
		s.address = net.JoinHostPort("", value)
	}
}

func WithAppName(value string) HttpOption {
	return func(s *HttpServer) {
		s.name = HttpServerPrefixName + value
	}
}

func WithPrintingRoutes(value bool) HttpOption {
	return func(s *HttpServer) {
		s.enablePrintingRoutes = value
	}
}
