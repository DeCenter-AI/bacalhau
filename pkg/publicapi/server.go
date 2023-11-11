package publicapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	echomiddelware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/time/rate"

	"github.com/bacalhau-project/bacalhau/docs"
	"github.com/bacalhau-project/bacalhau/pkg/config/types"
	"github.com/bacalhau-project/bacalhau/pkg/logger"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/middleware"
	"github.com/bacalhau-project/bacalhau/pkg/version"
)

const TimeoutMessage = "Server Timeout!"

type ServerParams struct {
	Router *echo.Echo
	HostID string

	Config types.ServerAPIConfig
}

// Server configures a node's public REST API.
type Server struct {
	Router *echo.Echo

	httpServer http.Server
	useTLS     bool
	Config     types.ServerAPIConfig
}

func NewAPIServer(params ServerParams) (*Server, error) {
	server := &Server{
		Router: params.Router,
		Config: params.Config,
	}

	// migrate old endpoints to new versioned ones
	migrations := map[string]string{
		"/peers":                      "/api/v1/peers",
		"/node_info":                  "/api/v1/node_info",
		"^/version":                   "/api/v1/version",
		"/healthz":                    "/api/v1/healthz",
		"/id":                         "/api/v1/id",
		"/livez":                      "/api/v1/livez",
		"/requester/list":             "/api/v1/requester/list",
		"/requester/nodes":            "/api/v1/requester/nodes",
		"/requester/states":           "/api/v1/requester/states",
		"/requester/results":          "/api/v1/requester/results",
		"/requester/events":           "/api/v1/requester/events",
		"/requester/submit":           "/api/v1/requester/submit",
		"/requester/cancel":           "/api/v1/requester/cancel",
		"/requester/debug":            "/api/v1/requester/debug",
		"/requester/logs":             "/api/v1/requester/logs",
		"/requester/websocket/events": "/api/v1/requester/websocket/events",
	}

	// set validator
	server.Router.Validator = NewCustomValidator()

	// enable debug mode to get clearer error messages
	// TODO: disable debug mode after we implement our own error handler
	server.Router.Debug = true

	// set middleware
	logLevel, err := zerolog.ParseLevel(params.Config.LogLevel)
	if err != nil {
		return nil, err
	}

	// base middleware before routing
	server.Router.Pre(
		echomiddelware.Rewrite(migrations),
	)

	// base middle after routing
	server.Router.Use(
		echomiddelware.TimeoutWithConfig(echomiddelware.TimeoutConfig{
			Timeout:      time.Duration(params.Config.RequestHandlerTimeout),
			ErrorMessage: TimeoutMessage,
			Skipper:      middleware.PathMatchSkipper(params.Config.SkippedTimeoutPaths),
		}),
		echomiddelware.RateLimiter(echomiddelware.NewRateLimiterMemoryStore(rate.Limit(params.Config.ThrottleLimit))),
		echomiddelware.RequestID(),
		middleware.RequestLogger(
			*log.Ctx(logger.ContextWithNodeIDLogger(context.Background(), params.HostID)),
			logLevel),
		middleware.Otel(),
		echomiddelware.BodyLimit(server.Config.MaxBytesToReadInBody),
		echomiddelware.Recover(),
	)

	if server.Config.EnableSwaggerUI {
		docs.SwaggerInfo.Version = version.Get().GitVersion
		server.Router.GET("/swagger/*", echo.WrapHandler(httpSwagger.WrapHandler))
	}

	var tlsConfig *tls.Config
	if params.Config.TLS.AutoCert != "" {
		log.Ctx(context.TODO()).Debug().Msgf("Setting up auto-cert for %s", params.Config.TLS.AutoCert)

		autoTLSManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(params.Config.TLS.AutoCertCachePath),
			HostPolicy: autocert.HostWhitelist(params.Config.TLS.AutoCert),
		}
		tlsConfig = &tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
			MinVersion:     tls.VersionTLS12,
		}

		server.useTLS = true
	} else {
		server.useTLS = params.Config.TLS.ServerCertificate != "" && params.Config.TLS.ServerKey != ""
	}

	server.Config.TLS.ServerCertificate = params.Config.TLS.ServerCertificate
	server.Config.TLS.ServerKey = params.Config.TLS.ServerKey

	server.httpServer = http.Server{
		Handler:           server.Router,
		ReadHeaderTimeout: time.Duration(server.Config.ReadHeaderTimeout),
		ReadTimeout:       time.Duration(server.Config.ReadTimeout),
		WriteTimeout:      time.Duration(server.Config.WriteTimeout),
		TLSConfig:         tlsConfig,
		BaseContext: func(l net.Listener) context.Context {
			return logger.ContextWithNodeIDLogger(context.Background(), params.HostID)
		},
	}

	return server, nil
}

// GetURI returns the HTTP URI that the server is listening on.
func (apiServer *Server) GetURI() *url.URL {
	interpolated := fmt.Sprintf("%s://%s:%d", apiServer.Config.Protocol, apiServer.Config.Host, apiServer.Config.Port)
	url, err := url.Parse(interpolated)
	if err != nil {
		panic(fmt.Errorf("callback url must parse: %s", interpolated))
	}
	return url
}

//	@title			Bacalhau API
//	@description	This page is the reference of the Bacalhau REST API. Project docs are available at https://docs.bacalhau.org/. Find more information about Bacalhau at https://github.com/bacalhau-project/bacalhau.
//	@contact.name	Bacalhau Team
//	@contact.url	https://github.com/bacalhau-project/bacalhau
//	@contact.email	team@bacalhau.org
//	@license.name	Apache 2.0
//	@license.url	https://github.com/bacalhau-project/bacalhau/blob/main/LICENSE
//	@host			bootstrap.production.bacalhau.org:1234
//	@BasePath		/
//	@schemes		http
//
// ListenAndServe listens for and serves HTTP requests against the API server.
//
//nolint:lll
func (apiServer *Server) ListenAndServe(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", apiServer.Config.Host, apiServer.Config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if apiServer.Config.Port == 0 {
		switch addr := listener.Addr().(type) {
		case *net.TCPAddr:
			apiServer.Config.Port = addr.Port
		default:
			return fmt.Errorf("unknown address %v", addr)
		}
	}

	log.Ctx(ctx).Debug().Msgf(
		"API server listening for host %s on %s...", apiServer.Config.Host, listener.Addr().String())

	go func() {
		var err error

		if apiServer.useTLS {
			err = apiServer.httpServer.ServeTLS(listener, apiServer.Config.TLS.ServerCertificate, apiServer.Config.TLS.ServerKey)
		} else {
			err = apiServer.httpServer.Serve(listener)
		}

		if err == http.ErrServerClosed {
			log.Ctx(ctx).Debug().Msgf(
				"API server closed for host %s on %s.", apiServer.Config.Host, apiServer.httpServer.Addr)
		} else if err != nil {
			log.Ctx(ctx).Err(err).Msg("Api server can't run. Cannot serve client requests!")
		}
	}()

	return nil
}

// Shutdown shuts down the http server
func (apiServer *Server) Shutdown(ctx context.Context) error {
	return apiServer.httpServer.Shutdown(ctx)
}
