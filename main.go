package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/ars0915/matching-system/config"
	"github.com/ars0915/matching-system/internal/tree"
	"github.com/ars0915/matching-system/router"
	"github.com/ars0915/matching-system/usecase"
	"github.com/ars0915/matching-system/util/log"
)

var (
	app        *cli.App
	drop       bool
	rollback   int
	configFile string

	// Version control.
	Version      = "No Version Provided"
	BuildDate    = ""
	GitCommitSha = ""
)

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "matching-system"
	app.Usage = "The RESTful service that provider web api"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "rollback",
			Value:       0,
			Destination: &rollback,
			Usage:       "rollback how many steps",
		},
		cli.StringFlag{
			Name:        "config, c",
			Value:       "",
			Destination: &configFile,
			Usage:       "Configuration file path",
		},
	}
	app.Action = func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			defer signal.Stop(quit)

			select {
			case <-ctx.Done():
			case <-quit:
				cancel()
			}
		}()

		// set default parameters.
		if err := config.InitConf(configFile); err != nil {
			logrus.Errorf("Load yaml config file error: '%v'", err)
			return err
		}

		logrus.WithFields(logrus.Fields{
			"logLevel": logrus.GetLevel(),
		}).Info("matching-system starting")

		log.SetLogLevel(config.Conf.Log.Level)

		boysTree := tree.NewPersonTree()
		girlsTree := tree.NewPersonTree()

		uHandler := usecase.InitHandler(boysTree, girlsTree)

		service := router.NewHandler(config.Conf, uHandler)

		if err := service.RunServer(ctx); err != nil {
			return err
		}

		return nil
	}
}

func main() {
	// Run the CLI app
	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Error("Service Run Fail")
	}
}
