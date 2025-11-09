package main

import (
	"net/url"
	"time"

	"github.com/ReanSn0w/gokit/pkg/app"
	"github.com/ReanSn0w/kincong/cmd/server/cr"
	"github.com/ReanSn0w/kincong/cmd/server/rest"
)

var (
	revision = "debug"
	opts     = struct {
		app.Debug

		Port    int    `long:"port" env:"PORT" default:"8080" description:"Port to listen on"`
		BaseURL string `long:"base-url" env:"BASE_URL" default:"http://localhost:8080" description:"Base URL for the server"`
	}{}
)

func main() {
	app := app.New("Network Data Resolver", revision, &opts)

	baseURL, err := url.Parse(opts.BaseURL)
	if err != nil {
		panic(err.Error())
	}

	restServer := rest.New(app.Log(), revision, baseURL, cr.New(app.Context()))
	app.Add(restServer.Stop)
	restServer.Start(app.CancelCause(), opts.Port)

	app.GracefulShutdown(time.Second * 10)
}
