package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var daemonStarted = false

func main() {
	var socketPath string

	flag.StringVar(&socketPath, "socket", "/run/guest-services/sqlcl-docker-extension.sock", "Unix domain socket to listen on")
	flag.Parse()

	os.RemoveAll(socketPath)

	logrus.New().Infof("Starting listening on %s", socketPath)
	router := echo.New()
	router.HideBanner = true

	startURL := ""

	ln, err := listen(socketPath)
	if err != nil {
		log.Fatal(err)
	}
	router.Listener = ln

	router.GET("/ready", ready)
	router.GET("/dark", dark)
	router.GET("/light", light)

	log.Fatal(router.Start(startURL))
}

func listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

// ready checks whether sqlcl is ready or not by querying localhost:9080.
func ready(ctx echo.Context) error {
	url := "http://localhost:7681/" // "sqlcl" is the name of the service defined in docker-compose.yml
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusOK, "false")

	}
	defer resp.Body.Close()

	return ctx.String(resp.StatusCode, "true")

	// return ctx.JSON(http.StatusOK, HTTPMessageBody{Message: "hello from HTTP"})

}

// dark starts ttyd in dark mode
func dark(ctx echo.Context) error {
	if !daemonStarted {
		cmd := exec.Command("/usr/bin/ttyd", "-u", "1000", "-g", "1000", "-t", "titleFixed='sqlcl'", "/bin/bash", "/home/sql.sh")
		cmd.Start()
		daemonStarted = true
	}
	return ctx.String(http.StatusOK, "true")

}

// light starts ttyd in dark mode
func light(ctx echo.Context) error {
	if !daemonStarted {
		cmd := exec.Command("/usr/bin/ttyd", "-u", "1000", "-g", "1000", "-t", "titleFixed='sqlcl'", "-t", "theme={'background': '#ffffff', 'foreground': '#2b2b2b', 'cursor': '#adadad', 'selection': '#ddb6fc'}", "/bin/bash", "/home/sql.sh")
		cmd.Start()
		daemonStarted = true
	}
	return ctx.String(http.StatusOK, "true")

}

type HTTPMessageBody struct {
	Message string `json:"message"`
	Body    string `json:"body,omitempty"`
}
