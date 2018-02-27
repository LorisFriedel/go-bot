package bot

import (
	"fmt"

	"github.com/LorisFriedel/go-bot/router"
	dgo "github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *dgo.Session
	Router  *router.Router
}

// New create a new instance of a discord bot that connect using the given token.
// The only possible additional argument is the bot name.
func New(token string) (*Bot, error) {
	b := &Bot{}

	err := b.initSession(token)
	if err != nil {
		return nil, err
	}

	err = b.initRouter()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bot) initSession(token string) error {
	session, err := dgo.New()
	if err != nil {
		return err
	}

	b.Session = session

	// Verify a Token was provided
	if token == "" {
		return fmt.Errorf("invalid Discord authentication token")
	}

	b.Session.Token = token

	// Open a websocket connection to Discord
	err = b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection to Discord (%v)", err)
	}

	return nil
}

func (b *Bot) initRouter() error {
	b.Router = router.New()

	// Register the Router OnMessageCreate handler that listens for and processes all messages received.
	b.Session.AddHandler(b.Router.OnMessageCreate)

	help, err := router.RouteBuilder.
		Prefix("go help").
		Description("Display this message").
		HandlerFunc(b.Router.BuiltinHelp).
		Build()

	if err != nil {
		return err
	}

	// Register the build-in help command.
	b.Router.AddRoute("help", help)

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}
