package router

import (
	"regexp"

	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lann/builder"
)

type IRoute interface {
	Matcher
	Handler
}

// HandlerFunc is the function signature required for a message route handler.
type HandlerFunc func(ds *discordgo.Session, mc *discordgo.MessageCreate, ch *discordgo.Channel)

// Handler contains the Handle method used as a callback for the router when a route match.
type Handler interface {
	Handle(ds *discordgo.Session, mc *discordgo.MessageCreate, ch *discordgo.Channel)
}

func (hf HandlerFunc) Handle(ds *discordgo.Session, mc *discordgo.MessageCreate, ch *discordgo.Channel) {
	hf(ds, mc, ch)
}

// HandlerFunc is the function signature required for a message route matcher.
// It is used to tell whether or not a message match a route.
type MatcherFunc func(msg string) bool

// Matcher contains the Match method used for the router to tell if a route match or not.
type Matcher interface {
	Match(msg string) bool
}

func (mf MatcherFunc) Match(msg string) bool {
	return mf(msg)
}

// Route holds information about a specific message route handler
type Route struct {
	prefix      string // match prefix that should trigger this route handler
	hasPrefix   bool
	contains    string // match contains that should trigger this route handler
	hasContains bool
	suffix      string // match suffix that should trigger this route handler
	hasSuffix   bool
	pattern     *regexp.Regexp // match regex pattern that should trigger this route handler
	hasPattern  bool
	description string  // short description of this route
	help        string  // detailed help string for this route
	handler     Handler // route handler function to call
	matcher     Matcher // route matcher function to call
	oneMatch    bool
}

// Handle is the callback called when this route is matched by a message.
func (r *Route) Handle(ds *discordgo.Session, mc *discordgo.MessageCreate, ch *discordgo.Channel) {
	if r.handler != nil {
		r.handler.Handle(ds, mc, ch)
	}
}

// Match check if the message fulfil the trigger matcher of this route. Optimized.
func (r *Route) Match(msg string) bool {
	if r.matcher != nil {
		return r.matcher.Match(msg)
	} else if r.oneMatch {
		return (r.hasPrefix && strings.HasPrefix(msg, r.prefix)) ||
			(r.hasContains && strings.Contains(msg, r.contains)) ||
			(r.hasSuffix && strings.HasSuffix(msg, r.suffix)) ||
			(r.hasPattern && r.pattern.MatchString(msg))
	} else {
		return (!r.hasPrefix || strings.HasPrefix(msg, r.prefix)) &&
			(!r.hasContains || strings.Contains(msg, r.contains)) &&
			(!r.hasSuffix || strings.HasSuffix(msg, r.suffix)) &&
			(!r.hasPattern || r.pattern.MatchString(msg))
	}
}

///////////////////////////////////////////
//////////////// Builder //////////////////
///////////////////////////////////////////

type rBuilder builder.Builder

// Builder used to simplify the creation of a route.
// Strict match by default (meaning prefix, suffix and pattern has to match for the method Match to return true)
var RouteBuilder = builder.Register(rBuilder{}, Route{}).(rBuilder)

// Condition that will trigger this route if it prefixes the message.
func (b rBuilder) Prefix(prefix string) rBuilder {
	return builder.Set(b, "prefix", prefix).(rBuilder)
}

// Condition that will trigger this route if the message contains it.
func (b rBuilder) Contains(contains string) rBuilder {
	return builder.Set(b, "contains", contains).(rBuilder)
}

// Condition that will trigger this route if it suffixes the message.
func (b rBuilder) Suffix(suffix string) rBuilder {
	return builder.Set(b, "suffix", suffix).(rBuilder)
}

// Regex expression that while trigger this route if it match the message.
func (b rBuilder) Pattern(pattern string) rBuilder {
	return builder.Set(b, "pattern", pattern).(rBuilder)
}

func (b rBuilder) Description(description string) rBuilder {
	return builder.Set(b, "description", description).(rBuilder)
}

func (b rBuilder) Help(help string) rBuilder {
	return builder.Set(b, "help", help).(rBuilder)
}

// Set the Handler to be the given function.
func (b rBuilder) HandlerFunc(handlerFunc HandlerFunc) rBuilder {
	return builder.Set(b, "handlerFunc", handlerFunc).(rBuilder)
}

// Set the Handler of this route. Override HandlerFunc if specified.
func (b rBuilder) Handler(handler Handler) rBuilder {
	return builder.Set(b, "handler", handler).(rBuilder)
}

// Set the Matcher to be the given function.
func (b rBuilder) MatcherFunc(matcherFunc MatcherFunc) rBuilder {
	return builder.Set(b, "matcherFunc", matcherFunc).(rBuilder)
}

// Set the Matcher of this route. Override MatcherFunc if specified.
func (b rBuilder) Matcher(matcher Matcher) rBuilder {
	return builder.Set(b, "matcher", matcher).(rBuilder)
}

// Tell the route to match if one or more matcher (prefix, suffix, pattern) matches.
func (b rBuilder) Soft() rBuilder {
	return builder.Set(b, "oneMatch", true).(rBuilder)
}

func (b rBuilder) Build() (*Route, error) {
	return newRoute(b)
}

func newRoute(rb rBuilder) (route *Route, err error) {
	var (
		prefix      string = ""
		contains    string = ""
		suffix      string = ""
		patternRe   *regexp.Regexp
		description string  = ""
		help        string  = ""
		handler     Handler = nil
		matcher     Matcher = nil
		oneMatch    bool    = false
	)

	if val, set := builder.Get(rb, "prefix"); set {
		prefix = val.(string)
	}

	if val, set := builder.Get(rb, "contains"); set {
		contains = val.(string)
	}

	if val, set := builder.Get(rb, "suffix"); set {
		suffix = val.(string)
	}

	if val, set := builder.Get(rb, "pattern"); set {
		pattern := val.(string)
		patternRe, err = regexp.Compile(pattern)
		if err != nil {
			return
		}
	}

	if val, set := builder.Get(rb, "description"); set {
		description = val.(string)
	}

	if val, set := builder.Get(rb, "help"); set {
		help = val.(string)
	}

	if val, set := builder.Get(rb, "handler"); set {
		handler = val.(Handler)
	} else if val, set := builder.Get(rb, "handlerFunc"); set {
		hf := val.(HandlerFunc)
		handler = hf
	}

	if val, set := builder.Get(rb, "matcher"); set {
		matcher = val.(Matcher)
	} else if val, set := builder.Get(rb, "matcherFunc"); set {
		mf := val.(MatcherFunc)
		matcher = mf
	}

	if val, set := builder.Get(rb, "oneMatch"); set {
		oneMatch = val.(bool)
	}

	return &Route{
		prefix:      prefix,
		hasPrefix:   prefix != "",
		contains:    contains,
		hasContains: contains != "",
		suffix:      suffix,
		hasSuffix:   suffix != "",
		pattern:     patternRe,
		hasPattern:  patternRe != nil,
		description: description,
		help:        help,
		handler:     handler,
		matcher:     matcher,
		oneMatch:    oneMatch,
	}, nil
}
