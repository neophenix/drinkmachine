package handlers

import (
	"github.com/neophenix/drinkmachine/internal/handlers/admin"
	"github.com/neophenix/drinkmachine/internal/handlers/api"
	"github.com/neophenix/drinkmachine/internal/handlers/ws"
	"net/http"
	"regexp"
)

// Route holds all our routing rules
type Route struct {
	Regex   *regexp.Regexp                               // a regex to compare the request path to
	Handler func(w http.ResponseWriter, r *http.Request) // a func pointer to call if the regex matches
}

// Routes is the array in the order we will attempt to match the route with the incoming url, first one wins
var Routes = []Route{
	{Regex: regexp.MustCompile("/$"), Handler: HomeHandler},
	{Regex: regexp.MustCompile("/api/pump$"), Handler: api.PumpListHandler},
	{Regex: regexp.MustCompile("/api/pump/[1-8]$"), Handler: api.PumpHandler},
	{Regex: regexp.MustCompile("/api/ingredient"), Handler: api.IngredientHandler},
	{Regex: regexp.MustCompile("/api/drink"), Handler: api.DrinkHandler},
	{Regex: regexp.MustCompile("/admin$"), Handler: admin.Handler},
	{Regex: regexp.MustCompile("/admin/pumps$"), Handler: admin.PumpHandler},
	{Regex: regexp.MustCompile("/admin/ingredients$"), Handler: admin.IngredientHandler},
	{Regex: regexp.MustCompile("/admin/drinks$"), Handler: admin.DrinkHandler},
	{Regex: regexp.MustCompile("/ws$"), Handler: ws.Handler},
}

// GetRouteHandler compares the path string to the route list and returns the handler pointer if found or nil
func GetRouteHandler(path string) func(w http.ResponseWriter, r *http.Request) {
	for _, route := range Routes {
		if route.Regex.MatchString(path) {
			return route.Handler
		}
	}

	return nil
}
