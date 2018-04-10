package bot

import (
	"fmt"

	"github.com/LorisFriedel/go-bot/router"
	dgo "github.com/LorisFriedel/discordgo"
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
	// Verify a Token was provided
	if token == "" {
		return fmt.Errorf("invalid empty Discord authentication token")
	}

	session, err := dgo.New("Bot " + token)
	if err != nil {
		return err
	}

	b.Session = session

	// Open a websocket connection to Discord
	err = session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection to Discord (%v)", err)
	}

	return nil
}

func (b *Bot) initRouter() error {
	b.Router = router.New()

	// Register the Router OnMessageCreate handler that listens for and processes all messages received.
	b.Session.AddHandler(b.Router.OnMessageCreate)

	// TODO allow to disable that or do something to easily add it if wanted
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
