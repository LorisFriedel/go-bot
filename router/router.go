// Package mux provides a simple Discord message route multiplexer that
// parses messages and then executes a matching registered handler, if found.
package router

import (
	"strings"

	dgo "github.com/LorisFriedel/discordgo"
	log "github.com/sirupsen/logrus"
)

// Router is the main struct for all router methods.
type Router struct {
	routes       map[string]IRoute
	defaultRoute IRoute
}

// New returns a new Discord message router
func New() *Router {
	return &Router{
		routes: make(map[string]IRoute),
	}
}

// Set the default route used if the message match nothing else.
func (r *Router) SetDefault(route IRoute) {
	r.defaultRoute = route
}

// Remove the default route. Messages that don't match any routes will be ignored.
func (r *Router) RemoveDefault() {
	r.defaultRoute = nil
}

// Add a new route to the router. Override any previous route with the same name.
func (r *Router) AddRoute(name string, route IRoute) {
	r.routes[name] = route
}

// Delete a route from the router. Does nothing if no such route exists.
func (r *Router) DeleteRoute(name string) {
	delete(r.routes, name)
}

// matchRoutes attempts to find the best route match for a givin message.
func (r *Router) matchRoutes(msg string) []IRoute {
	routes := make([]IRoute, 0, 1)

	for _, r := range r.routes {
		if r.Match(msg) {
			routes = append(routes, r)
		}
	}

	if len(routes) == 0 && r.defaultRoute != nil {
		routes = append(routes, r.defaultRoute)
	}

	return routes
}

// OnMessageCreate is a DiscordGo Event Handler function.  This must be
// registered using the DiscordGo.Session.AddHandler function.  This function
// will receive all Discord messages and parse them for matches to registered
// routes.
func (r *Router) OnMessageCreate(ds *dgo.Session, mc *dgo.MessageCreate) {

	// Ignore all messages created by the Bot account itself
	if mc.Author.ID == ds.State.User.ID {
		return
	}

	// Fetch the channel for this Message
	ch, err := ds.State.Channel(mc.ChannelID)
	if err != nil {
		// Try fetching via REST API
		ch, err = ds.Channel(mc.ChannelID)
		if err != nil {
			log.Errorf("unable to fetch Channel for Message")
			return
		}
		// Attempt to add this channel into our State
		err = ds.State.ChannelAdd(ch)
		if err != nil {
			log.Errorf("error updating State with Channel")
		}
	}

	// Match message content against every routes
	routes := r.matchRoutes(strings.TrimSpace(mc.Content))

	// Execute routes that matched
	for _, r := range routes {
		go r.Handle(ds, mc, ch)
	}
}

// Help function provides a build in "help" command that will display a list
// of all registered routes (commands). To use this function it must first be
// registered with the Router.AddRoute function.
func (r *Router) BuiltinHelp(ds *dgo.Session, mc *dgo.MessageCreate, ch *dgo.Channel) {

	// TODO
	// ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
