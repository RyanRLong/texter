package main

import (
	"log"
	"net/smtp"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/jordan-wright/email"
	"github.com/urfave/cli"
)

// Config file structure
type Config struct {
	From     string
	To       []string
	Username string
	Password string
	Server   string
	Port     int
}

// external config
var conf Config

// Gets the configuration from a config.toml file
func getConfig() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
}

// initiate cli app
var app = cli.NewApp()

// placeholder for subject
var subject string

// verbosity flag
var verbose bool

// attaches info to the cli app
func initInfo() {
	app.Name = "texter"
	app.Usage = "send a text"
	app.Author = "Ryan Long"
	app.Version = "1.0.0"
}

// attaches flags to the cli app
func initFlags() {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "subject, s",
			Value:       "",
			Usage:       "`SUBJECT` for the SMS, defaults to none",
			Destination: &subject,
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose output",
		},
	}
}

// attaches action to the cli app
func getAction() {
	verbose = false

	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 || c.Args().Get(0) == "" {
			return cli.NewExitError("A valid message is required", 1)
		}

		if c.Bool("verbose") {
			verbose = true
		}

		message := []byte(c.Args().Get(0))

		if verbose {
			log.Printf("Preparing to send SMS with message: %s and subject: %s", message, subject)
		}
		err := send(message, subject)
		if err != nil {
			log.Fatal(err)
		}
		if verbose {
			log.Print("Message sent successfully")
		}
		return nil
	}
}

// sends the message with an options subject.  message is required.
func send(message []byte, subject string) error {
	if verbose {
		log.Printf("Attempting to send SMS")
	}
	uri := conf.Server + ":" + strconv.Itoa(conf.Port)
	e := email.NewEmail()
	e.From = conf.From
	e.To = conf.To
	e.Subject = subject
	e.Text = message
	if verbose {
		log.Printf("Parameters are: %v", e)
	}
	// e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	return e.Send(uri, smtp.PlainAuth("", conf.Username, conf.Password, conf.Server))
}

// Main point of execution
func main() {
	getConfig()
	initInfo()
	initFlags()
	getAction()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
