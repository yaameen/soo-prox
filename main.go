package main

import (
	"log"
	"os"
	"sooprox/cmd"
	"sooprox/types"
	"sooprox/utils"

	"github.com/urfave/cli/v2"
)

func main() {
	flags := []cli.Flag{
		// config
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
			Value:   "config.yaml",
			EnvVars: []string{"CONFIG_FILE"},
		},
		// host
		&cli.StringFlag{
			Name:    "host",
			Aliases: []string{"H"},
			Usage:   "Listen on `HOST`",
			EnvVars: []string{"HOST"},
		},
		// port
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Listen on `PORT`",
			EnvVars: []string{"PORT"},
		},
		// proxies a string with format prefix:target
		&cli.StringSliceFlag{
			Name:    "proxies",
			Aliases: []string{"P"},
			Value:   cli.NewStringSlice(),
			Usage:   "Proxy `PREFIX::TARGET`",
		},
		// is secure
		&cli.BoolFlag{
			Name:    "secure",
			Aliases: []string{"s"},
			Value:   false,
			Usage:   "Use TLS",
		},
	}
	app := &cli.App{
		Name:    "SooProx",
		Usage:   "A simple proxy server",
		Version: cmd.Version,
		Authors: []*cli.Author{
			{
				Name:  "Yameen Mohamed",
				Email: "yaamynu@gmail.com",
			},
		},
		Action: func(ctx *cli.Context) error {

			isCliMode := false

			host := ctx.String("host")
			port := ctx.Int("port")
			config := ctx.String("config")

			if host == "" || port == 0 {
				isCliMode = true

			}
			proxies := ctx.StringSlice("proxies")

			fConfig := types.Config{
				Secure: ctx.Bool("secure"),
			}

			if len(proxies) > 0 {
				fConfig.Proxies = []types.Proxy{}
				for _, v := range proxies {
					fConfig.Proxies = append(fConfig.Proxies, utils.ParseProxy(v))
				}
			}
			if _, err := os.Stat(config); err == nil {
				log.Printf("Loading config from %s", config)
				configFromFile := utils.ReadConfig(config)
				if host == "" {
					host = configFromFile.Host
				}
				if port == 0 {
					port = configFromFile.Port
				}
				fConfig.ConfigFile = ctx.String("config")
				fConfig.Proxies = append(fConfig.Proxies, configFromFile.Proxies...)
			}

			if host == "" {
				host = "localhost"
			}
			if port == 0 {
				port = 8080
			}
			fConfig.Host = host
			fConfig.Port = port

			if port == 443 {
				fConfig.Secure = true
			}

			cmd.Init(fConfig, isCliMode)

			return nil
		},
		Flags: flags,
		Commands: []*cli.Command{
			{
				Name:    "ca-trust",
				Aliases: []string{"t"},
				Usage:   "Trust the CA certificate",
				Action: func(_ *cli.Context) error {
					utils.TrustCA()
					return nil
				},
			},
			{
				Name:    "ca-gen",
				Aliases: []string{"g"},
				Usage:   "Generate a CA certificate",
				Action: func(_ *cli.Context) error {
					utils.GenerateCA()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
