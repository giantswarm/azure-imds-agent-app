package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/giantswarm/microerror"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/giantswarm/azure-imds-agent-app/flag"
	httpclusterid "github.com/giantswarm/azure-imds-agent-app/server/context/http/clusterid"
	httpheader "github.com/giantswarm/azure-imds-agent-app/server/context/http/header"
	"github.com/giantswarm/azure-imds-agent-app/server/context/http/instanceid"
	httpmethod "github.com/giantswarm/azure-imds-agent-app/server/context/http/method"
	httpnodepoolid "github.com/giantswarm/azure-imds-agent-app/server/context/http/nodepoolid"
	httporganizationid "github.com/giantswarm/azure-imds-agent-app/server/context/http/organizationid"
	"github.com/giantswarm/azure-imds-agent-app/server/context/http/resourcegroupname"
	"github.com/giantswarm/azure-imds-agent-app/server/context/http/subscriptionid"
	httpurl "github.com/giantswarm/azure-imds-agent-app/server/context/http/url"
	"github.com/giantswarm/azure-imds-agent-app/server/context/http/vmssname"
	"github.com/giantswarm/azure-imds-agent-app/server/endpoint"
	instancegetter "github.com/giantswarm/azure-imds-agent-app/server/endpoint/imds/instance/getter"
	"github.com/giantswarm/azure-imds-agent-app/server/middleware"
	"github.com/giantswarm/azure-imds-agent-app/service"
	"github.com/giantswarm/azure-imds-agent-app/service/option"
)

// Config represents the configuration used to create a new server object.
type Config struct {
	Flag    *flag.Flag
	Logger  micrologger.Logger
	Service *service.Service
	Viper   *viper.Viper

	ProjectName string
}

// Server is the star of the show.
type Server struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	bootOnce     sync.Once
	config       microserver.Config
	shutdownOnce sync.Once
}

// New creates a new configured server object.
func New(config Config) (*Server, error) {
	var err error

	var middlewareCollection *middleware.Middleware
	{
		c := middleware.Config{
			Flag:    config.Flag,
			Logger:  config.Logger,
			Service: config.Service,
			Viper:   config.Viper,
		}

		middlewareCollection, err = middleware.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoints *endpoint.Endpoint
	{
		c := endpoint.DefaultConfig()
		c.Flag = config.Flag
		c.Logger = config.Logger
		c.Middleware = middlewareCollection
		c.Service = config.Service
		c.Viper = config.Viper

		endpoints, err = endpoint.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// corsHeaders, err := config.Service.Option.Options(context.Background(), option.Request{
	// 	AllowedOrigin: config.Viper.GetString(config.Flag.Server.AllowedOrigin),
	// })
	// if err != nil {
	// 	return nil, microerror.Mask(err)
	// }

	handlerWrapper, err := createHandlerWrapper([]*option.Response{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	s := &Server{
		// Dependencies.
		logger: config.Logger,

		// Internals.
		bootOnce: sync.Once{},
		config: microserver.Config{
			Logger:      config.Logger,
			ServiceName: config.ProjectName,
			Viper:       config.Viper,

			Endpoints: []microserver.Endpoint{
				endpoints.IMDS.Instance.Getter,
			},
			ErrorEncoder:   errorEncoder,
			HandlerWrapper: handlerWrapper,
			RequestFuncs:   newRequestFuncs(),
		},
		shutdownOnce: sync.Once{},
	}

	return s, nil
}

// Boot boots the server.
func (s *Server) Boot() {
	s.bootOnce.Do(func() {
		// Here goes your custom boot logic for your server/endpoint/middleware, if
		// any.
	})
}

// Config returns the configuration of this server.
func (s *Server) Config() microserver.Config {
	return s.config
}

// Shutdown tells the server to shut itself down.
func (s *Server) Shutdown() {
	s.shutdownOnce.Do(func() {
		// Here goes your custom shutdown logic for your server/endpoint/middleware,
		// if any.
	})
}

func createHandlerWrapper(headers []*option.Response) (func(http.Handler) http.Handler, error) {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, h := range headers {
				w.Header().Set(h.Key, h.Value)
			}
			w.Header().Set("connection", "keep-alive")

			if r.Method == "OPTIONS" {
				return
			}

			h.ServeHTTP(w, r)
		})
	}, nil
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	rErr := err.(microserver.ResponseError)
	uErr := rErr.Underlying()

	if instancegetter.IsNotFound(uErr) { // check for instance metadata not found error
		rErr.SetCode(microserver.CodeResourceNotFound)
		rErr.SetMessage(fmt.Sprintf("The instance metadata could not be found. (%s)", uErr.Error()))
		w.WriteHeader(http.StatusNotFound)
	} else {
		rErr.SetCode(microserver.CodeInternalError)
		rErr.SetMessage("An unexpected error occurred. Sorry for the inconvenience.")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// This writes the error response body to the stream.
	encodingErr := json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    rErr.Code(),
		"message": rErr.Message(),
	})
	if encodingErr != nil {
		panic(encodingErr)
	}
}

func newRequestFuncs() []kithttp.RequestFunc {
	return []kithttp.RequestFunc{
		// This request function puts the cluster ID URL parameter into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return httpclusterid.NewContext(ctx, mux.Vars(r)["cluster_id"])
		},
		// This request function puts the organization ID URL parameter into the
		// request context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return httporganizationid.NewContext(ctx, mux.Vars(r)["organization_id"])
		},
		// This request function puts the node pool ID URL parameter into the
		// request context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return httpnodepoolid.NewContext(ctx, mux.Vars(r)["nodepool_id"])
		},
		// This request function puts the subscription ID URL parameter into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return subscriptionid.NewContext(ctx, mux.Vars(r)["subscription_id"])
		},
		// This request function puts the resource group name URL parameter into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return resourcegroupname.NewContext(ctx, mux.Vars(r)["resource_group_name"])
		},
		// This request function puts the VMSS name URL parameter into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return vmssname.NewContext(ctx, mux.Vars(r)["vmss_name"])
		},
		// This request function puts the instance ID URL parameter into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return instanceid.NewContext(ctx, mux.Vars(r)["instance_id"])
		},
		// This request function puts the HTTP Authorization header into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return httpheader.NewContext(ctx, r.Header.Get("Authorization"))
		},
		// This request function puts the HTTP method into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return httpmethod.NewContext(ctx, r.Method)
		},
		// This request function puts the HTTP URL into the request
		// context, if any.
		func(ctx context.Context, r *http.Request) context.Context {
			return httpurl.NewContext(ctx, r.URL.String())
		},
	}
}
