package main

import (
	"context"
	alertHttp "github.com/ns3777k/alertmanager-webhook-space/http"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ns3777k/alertmanager-webhook-space/pkg/space"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"gopkg.in/alecthomas/kingpin.v2"
)

func makeSignalShutdownChan() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	return c
}

type Configuration struct {
	ListenAddr string
	ChannelID  string
}

func main() {
	godotenv.Load() //nolint:golint,errcheck

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	configuration := new(Configuration)
	spaceClientSettings := new(space.ClientSettings)

	app := kingpin.New(filepath.Base(os.Args[0]), "alertmanager-webhook-space")
	app.HelpFlag.Short('h')

	app.Flag("listen", "Address to listen on").
		Short('l').
		Default("0.0.0.0:9091").
		Envar("LISTEN_ADDR").
		StringVar(&configuration.ListenAddr)

	app.Flag("space-channel-id", "Channel id").
		Envar("CHANNEL_ID").
		Required().
		StringVar(&configuration.ChannelID)

	app.Flag("space-base-url", "Base url like https://mycompany.jetbrains.space").
		Envar("BASE_URL").
		Required().
		StringVar(&spaceClientSettings.BaseURL)

	app.Flag("space-client-id", "Application client id").
		Envar("CLIENT_ID").
		Required().
		StringVar(&spaceClientSettings.ClientID)

	app.Flag("space-client-secret", "Application client secret").
		Envar("CLIENT_SECRET").
		Required().
		StringVar(&spaceClientSettings.ClientSecret)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.Usage(os.Args[1:])
		logger.Fatal(err)
	}

	shutdownCh := make(chan struct{})
	spaceClient := space.NewClient(spaceClientSettings)
	h := alertHttp.NewHandler(spaceClient, configuration.ChannelID)
	server := &http.Server{
		ReadTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		Addr:              configuration.ListenAddr,
		Handler:           h.Router,
	}

	go func() {
		<-makeSignalShutdownChan()

		graceCtx, graceCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer graceCancel()

		if err := server.Shutdown(graceCtx); err != nil {
			logger.WithError(err).Error("error while shutting down the server")
		}

		shutdownCh <- struct{}{}
	}()

	logger.Info("start listening on " + configuration.ListenAddr)
	err := server.ListenAndServe()
	<-shutdownCh

	logger.Info("shutting down the server")

	if err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}
