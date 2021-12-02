package web

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/chuckha/geogame6/internal/usecases"
)

type adapterUsecases interface {
	GetCountryData(ctx context.Context, in *usecases.GetCountryDataInput) (*usecases.GetCountryDataOutput, error)
}

type HTTPAdapter struct {
	usecases  adapterUsecases
	presenter *presenter
}

func NewHTTPAdapter(usecases adapterUsecases) *HTTPAdapter {
	return &HTTPAdapter{
		usecases:  usecases,
		presenter: &presenter{},
	}
}

func NewServer(adapter *HTTPAdapter, mux *http.ServeMux, addr string) *http.Server {
	mux.HandleFunc("/get-country-data", adapter.getCountryData)
	return &http.Server{
		Handler: mux,
		Addr:    addr,
	}
}

func (h *HTTPAdapter) getCountryData(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	in := &usecases.GetCountryDataInput{Name: country}
	out, err := h.usecases.GetCountryData(r.Context(), in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		http.Error(w, "failed to get country data", http.StatusInternalServerError)
		return
	}
	if err := h.presenter.presentGetCountryData(w, out); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		http.Error(w, "failed to present data", http.StatusInternalServerError)
		return
	}
}
