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
	// Parse arguments
	flag.Parse()
	argsCli := parseCli()
	argsEnvVar := parsEnvVar()
	args := merge(argsEnvVar, argsCli)

	// Set up
	gobot, err := bot.New(args.token)
	if err != nil {
		log.Errorln(err)
		return
	}

	// TODO add my amazing routes like:
	// gobot.Router.AddRoute(
	// 				"amazingFunctionality",
	//				router.RouteBuilder.Contains("fromage").HandlerFunc(...).Build())

	log.Infoln("The Bot is running. Press CTRL-C to exit.")
	waitToBeMurdered()

	// Clean up
	err = gobot.Stop()
	if err != nil {
		log.Errorln(err)
		return
	}
}

// Wait for a CTRL-C
func waitToBeMurdered() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

///////////////////////////////////////////
//////////////// Parsing //////////////////
///////////////////////////////////////////

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
// (the last one can override the first one but not the other way round)
func merge(argsList ...*Arguments) *Arguments {
	result := &Arguments{}
	for _, args := range argsList {
		if args.token != "" {
			result.token = args.token
		}
	}
	return result
}
