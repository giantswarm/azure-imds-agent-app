package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/azure-imds-agent-app/flag"
	"github.com/giantswarm/azure-imds-agent-app/pkg/project"
	"github.com/giantswarm/azure-imds-agent-app/server"
	"github.com/giantswarm/azure-imds-agent-app/service"
)

var (
	f *flag.Flag = flag.New()
)

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		c := micrologger.Config{}

		newLogger, err = micrologger.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		var newService *service.Service
		{
			httpClient := http.DefaultClient
			c := service.Config{
				Logger: newLogger,

				Flag:  f,
				Viper: v,

				Description: project.Description(),
				GitCommit:   project.GitSHA(),
				Name:        project.Name(),
				Source:      project.Source(),

				HttpClient: httpClient,
			}

			newService, err = service.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
		}

		var newServer microserver.Server
		{
			c := server.Config{
				Flag:    f,
				Logger:  newLogger,
				Service: newService,
				Viper:   v,

				ProjectName: project.Name(),
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        newLogger,
			ServerFactory: newServerFactory,

			Description: project.Description(),
			GitCommit:   project.GitSHA(),
			Name:        project.Name(),
			Source:      project.Source(),
			Version:     project.Version(),
		}

		newCommand, err = command.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()
	daemonCommand.PersistentFlags().String(f.Server.AllowedOrigin, "*", "The value for the Access-Control-Allow-Origin header.")
	daemonCommand.PersistentFlags().Duration(f.Middleware.RateLimiter.DefaultPeriod, 5*time.Minute, "Default time period for the rate limiter middleware.")
	daemonCommand.PersistentFlags().Int64(f.Middleware.RateLimiter.DefaultLimit, 5, "Default limit of requests for the rate limiter middleware.")

	err = newCommand.CobraCommand().Execute()
	if err != nil {
		panic(err)
	}
}
