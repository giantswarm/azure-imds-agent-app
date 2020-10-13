package getter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	jsoniter "github.com/json-iterator/go"

	"github.com/giantswarm/azure-imds-agent-app/server/context/http/instanceid"
	"github.com/giantswarm/azure-imds-agent-app/server/context/http/resourcegroupname"
	"github.com/giantswarm/azure-imds-agent-app/server/context/http/subscriptionid"
	"github.com/giantswarm/azure-imds-agent-app/server/context/http/vmssname"
	"github.com/giantswarm/azure-imds-agent-app/server/middleware"
	"github.com/giantswarm/azure-imds-agent-app/service/imds/instance"
)

const (
	// Method is the HTTP method this endpoint is registered for.
	Method = "GET"
	// Name identifies the endpoint. It is aligned to the package path.
	Name = "imds/instance/getter/"
	// Path is the HTTP request path this endpoint is registered for.
	Path = "/imds/instance/{subscription_id}/{resource_group_name}/{vmss_name}/{instance_id}/"
)

type Config struct {
	Logger          micrologger.Logger
	Middleware      *middleware.Middleware
	InstanceService *instance.Service
}

type Endpoint struct {
	logger          micrologger.Logger
	middleware      *middleware.Middleware
	instanceService *instance.Service
}

// New creates a new configured endpoint object.
func New(config Config) (*Endpoint, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Middleware must not be empty", config)
	}
	if config.InstanceService == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.InstanceService must not be empty", config)
	}

	e := &Endpoint{
		logger:          config.Logger,
		middleware:      config.Middleware,
		instanceService: config.InstanceService,
	}

	return e, nil
}

// Decoder decodes the incoming request.
func (e *Endpoint) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		subscriptionID := subscriptionid.FromContext(ctx)
		resourceGroupName := resourcegroupname.FromContext(ctx)
		vmssName := vmssname.FromContext(ctx)
		instanceID := instanceid.FromContext(ctx)

		request := Request{
			PathParams: PathParams{
				SubscriptionID:    subscriptionID,
				ResourceGroupName: resourceGroupName,
				VMSSName:          vmssName,
				InstanceID:        instanceID,
			},
		}

		return request, nil
	}
}

// Encoder encodes the outgoing response.
func (e *Endpoint) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		err := jsoniter.NewEncoder(w).Encode(response)
		if err != nil {
			return microerror.Mask(err)
		}
		return nil
	}
}

// Endpoint is where the actual request gets sent off to a service package that does some business logic.
func (e *Endpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		endpointRequest := request.(Request)

		e.logger.LogCtx(ctx,
			"level", "debug",
			"message", fmt.Sprintf(
				"getting instance metadata for subscription %s, resource group %s, VMSS %s, instance %s",
				endpointRequest.PathParams.SubscriptionID,
				endpointRequest.PathParams.ResourceGroupName,
				endpointRequest.PathParams.VMSSName,
				endpointRequest.PathParams.InstanceID))

		instanceMetadata, err := e.instanceService.GetMetadata(ctx)
		if err != nil {
			e.logger.LogCtx(ctx,
				"level", "error",
				"message", fmt.Sprintf(
					"error while getting instance metadata for subscription %s, resource group %s, VMSS %s, instance %s",
					endpointRequest.PathParams.SubscriptionID,
					endpointRequest.PathParams.ResourceGroupName,
					endpointRequest.PathParams.VMSSName,
					endpointRequest.PathParams.InstanceID))
			return nil, microerror.Mask(err)
		}

		return instanceMetadata, nil
	}
}

// Method returns what HTTP verb this endpoint accepts. These are defined as consts above.
func (e *Endpoint) Method() string {
	return Method
}

// Middlewares returns a slice of the middlewares used in this endpoint.
func (e *Endpoint) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

// Name returns a name for this endpoint. These are defined as consts above.
func (e *Endpoint) Name() string {
	return Name
}

// Path returns thet path that this endpoint is listening on. These are defined as consts above.
func (e *Endpoint) Path() string {
	return Path
}
