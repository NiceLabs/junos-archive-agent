package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/google/go-github/github"
	"goftp.io/server"
	"golang.org/x/oauth2"
)

var conf = &Configuration{
	FTP: &FTPConfigure{
		Port:     21,
		Username: "anonymous",
	},
	GitHub: &GitHubConfigure{
		Prefix: "Juniper",
	},
}

func init() {
	cfgPath := "./configure.json"
	if data, err := ioutil.ReadFile(cfgPath); err != nil {
		panic(err)
	} else if err = json.Unmarshal(data, conf); err != nil {
		panic(err)
	}
}

func main() {
	factory := &junosDriver{
		GitHubConfigure: conf.GitHub,
		client: github.NewClient(
			oauth2.NewClient(
				context.Background(),
				oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.GitHub.Token}),
			),
		),
	}
	serve := server.NewServer(&server.ServerOpts{
		Name:           "JunOS Archive",
		WelcomeMessage: "JunOS FTP Archive Server",
		Port:           conf.FTP.Port,
		Hostname:       conf.FTP.Host,
		Auth:           makeAuth(),
		Factory:        factory,
	})
	if err := serve.ListenAndServe(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func makeAuth() server.Auth {
	name := conf.FTP.Username
	password := conf.FTP.Password
	if (name == "anonymous" || name == "") && password == "" {
		return &AnonymousAuth{}
	}
	return &server.SimpleAuth{Name: name, Password: password}
}
