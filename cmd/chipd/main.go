package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"s.mcquay.me/dm/chipmunk"

	"github.com/gorilla/context"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

type Config struct {
	Host         string
	Port         int
	DBHost       string
	DBName       string
	ClientID     string
	ClientSecret string
	CookieSecret string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var host, clientId, clientSecret, cookieSecret string
	var port int

	var run = &cobra.Command{
		Use:   "run",
		Short: "run command",
		Long:  `run chipd with given options`,
		Run: func(cmd *cobra.Command, args []string) {
			config := &Config{
				DBHost: "localhost",
				DBName: "dmcquay",
			}
			err := envconfig.Process("chipd", config)
			if err != nil {
				log.Fatal(err)
			}
			if host != "" {
				config.Host = host
			}
			if port != -1 {
				config.Port = port
			} else {
				if config.Port == 0 {
					config.Port = 8080
				}
			}
			if clientId != "" {
				config.ClientID = clientId
			}
			if clientSecret != "" {
				config.ClientSecret = clientSecret
			}
			if cookieSecret != "" {
				config.CookieSecret = cookieSecret
			} else {
				if config.CookieSecret == "" {
					config.CookieSecret = strconv.Itoa(rand.Int())
				}
			}
			if config.ClientID == "" || config.ClientSecret == "" {
				log.Println("Please provide both a clientId and a clientSecret")
				os.Exit(1)
			}
			log.Printf("%+v", config)

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, os.Interrupt, os.Kill)
			go func() {
				s := <-sigs
				log.Printf("signal: %+v", s)
				os.Exit(1)
			}()

			sm := http.NewServeMux()
			_, err = chipmunk.NewServer(
				sm,
				config.ClientID,
				config.ClientSecret,
				config.CookieSecret,
				config.DBHost,
				config.DBName,
				"",
			)
			if err != nil {
				log.Fatalf("problem initializing Chipd server: %+v", err)
			}

			hostname := "localhost"
			if config.Host == "" {
				hostname, err = os.Hostname()
				if err != nil {
					log.Printf("problem getting hostname:", err)
				}
			}
			log.Printf("serving at: http://%s:%d/", hostname, config.Port)

			addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
			err = http.ListenAndServe(addr, context.ClearHandler(sm))
			if err != nil {
				log.Printf("%+v", err)
				os.Exit(1)
			}
		},
	}

	run.Flags().StringVarP(
		&host,
		"host",
		"n",
		"",
		"hostname",
	)
	run.Flags().IntVarP(
		&port,
		"port",
		"p",
		-1,
		"port",
	)
	run.Flags().StringVarP(
		&clientId,
		"clientId",
		"i",
		"",
		"github oauth clientId",
	)
	run.Flags().StringVarP(
		&clientSecret,
		"clientSecret",
		"s",
		"",
		"github oauth clientSecret",
	)
	run.Flags().StringVarP(
		&cookieSecret,
		"cookieSecret",
		"c",
		"",
		"cookieSecret",
	)

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(run)
	rootCmd.Execute()
}
