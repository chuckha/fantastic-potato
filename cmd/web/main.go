package main

import (
	"net/http"

	"github.com/chuckha/geogame6/internal/adapters/web"
	"github.com/chuckha/geogame6/internal/infrastructure/encdec/json"
	"github.com/chuckha/geogame6/internal/infrastructure/gatherers/geojson/local"
	"github.com/chuckha/geogame6/internal/usecases"
)

func main() {
	statesFile := "./data/states.json"
	out := "./data"
	gatherer := local.NewGatherer(statesFile, out)
	encdec := &json.EncDec{}
	countryData := usecases.NewCountryData(gatherer, encdec)
	adapter := web.NewHTTPAdapter(countryData)
	mux := http.NewServeMux()
	server := web.NewServer(adapter, mux, ":8888")
	server.ListenAndServe()
}

// layers:
// infra: db (definitely need persistence)
// <input> -> adapter converts input to usecase input
// call usecase
// presenter takes usecase output and returns the output -> <output>

// httpadapter
// accepts requests
// httpAdapter := NewHTTPAdapter()
// usecases := NewUseCases()
//
// webpresenter
// returns html,js,css
