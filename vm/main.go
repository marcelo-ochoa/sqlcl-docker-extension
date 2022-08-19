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

var ttyd TTYD

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

	ttyd = TTYD{}

	router.GET("/ready", ready)
	router.POST("/dark", dark)
	router.POST("/light", light)

	log.Fatal(router.Start(startURL))
}

func listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

// ready checks whether sqlcl is ready or not by querying localhost:9080.
func ready(ctx echo.Context) error {
	if ttyd.IsRunning() {
		return ctx.String(http.StatusOK, "true")
	}

	return ctx.String(http.StatusServiceUnavailable, "false")
}

// dark starts ttyd in dark mode
func dark(ctx echo.Context) error {
	if err := ttyd.Start(DarkTheme); err != nil {
		log.Printf("failed to start ttyd with dark mode: %s\n", err)

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Println("started ttyd with dark mode")

	return ctx.String(http.StatusOK, "true")
}

// light starts ttyd in dark mode
func light(ctx echo.Context) error {
	if err := ttyd.Start(LightTheme); err != nil {
		log.Printf("failed to start ttyd with light mode: %s\n", err)

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Println("started ttyd with light mode")

	return ctx.String(http.StatusOK, "true")
}

type HTTPMessageBody struct {
	Message string `json:"message"`
	Body    string `json:"body,omitempty"`
}

type Theme string

const (
	DarkTheme  Theme = "dark"
	LightTheme Theme = "light"
)

type TTYD struct {
	process *os.Process
}

func (t *TTYD) Start(theme Theme) error {
	if t.IsStarted() {
		if err := t.Stop(); err != nil {
			log.Printf("failed to stop ttyd: %s\n", err)
		}
	}

	args := []string{"-u", "1000", "-g", "1000", "-t", "titleFixed='sqlcl'"}
	if theme == LightTheme {
		args = append(args, "-t", "theme={'background': '#ffffff', 'foreground': '#2b2b2b', 'cursor': '#adadad', 'selection': '#ddb6fc'}")
	}

	args = append(args, "/bin/bash", "/home/sql.sh")

	cmd := exec.Command("/usr/bin/ttyd", args...)
	if err := cmd.Start(); err != nil {
		return err
	}

	t.process = cmd.Process

	return nil
}

func (t *TTYD) Stop() error {
	if !t.IsStarted() {
		return nil
	}

	if err := t.process.Kill(); err != nil {
		return err
	}

	return nil
}

func (t TTYD) IsStarted() bool {
	return t.process != nil
}

func (t *TTYD) IsRunning() bool {
	if !t.IsStarted() {
		return false
	}

	url := "http://localhost:7681/" // "sqlcl" is the name of the service defined in docker-compose.yml
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false

	}

	return resp.StatusCode == http.StatusOK
}
