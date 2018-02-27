package main

import (
	"flag"

	"os"
	"os/signal"
	"syscall"

	"github.com/LorisFriedel/go-bot/bot"
	log "github.com/sirupsen/logrus"
)

type Arguments struct {
	token string
}

var argToken string

func init() {
	flag.StringVar(&argToken, "t", "", "Discord Authentication Token")
}

func main() {
	// Parse arguments from all sources
	flag.Parse()
	args := merge(parsEnvVar(), parseCli())

	// Set up
	gobot, err := bot.New(args.token)
	if err != nil {
		log.Errorln(err)
		return
	}

	// TODO add my amazing simple routes like:
	// myRoute, err := router.RouteBuilder.Contains("fromage").HandlerFunc(...).Build()
	// gobot.Router.AddRoute("amazingFunctionality", myRoute)
	// TODO TEXT ROUTER

	// TODO add voice route using a special router?
	// TODO VOICE ROUTER
	// todo + possibility to overwrite or add new voice handler

	log.Infoln("The Bot is running. Press CTRL-C to exit.")
	waitSIGTERM()

	// Clean up
	err = gobot.Stop()
	if err != nil {
		log.Errorln(err)
		return
	}
}

// Wait for a CTRL-C
func waitSIGTERM() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

/////////////////////////////////////////////////////
//////////////// Parsing arguments //////////////////
/////////////////////////////////////////////////////

func parsEnvVar() *Arguments {
	return &Arguments{
		token: os.Getenv("D_TOKEN"),
	}
}

func parseCli() *Arguments {
	return &Arguments{
		token: argToken,
	}
}

// merge aggregate arguments from every given sources. Override are made regarding the method arguments order.
// (the last one may override the first one but not the other way round)
func merge(argsList ...*Arguments) *Arguments {
	result := &Arguments{}
	for _, args := range argsList {
		if args.token != "" {
			result.token = args.token
		}
	}
	return result
}
