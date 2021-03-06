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

func NewAPIHandler(adapter *HTTPAdapter) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/get-country-data", adapter.getCountryData)
	return mux
}

func (h *HTTPAdapter) getCountryData(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	targetLang := r.URL.Query().Get("target_language")
	if targetLang == "" {
		targetLang = "japanese"
	}
	in := &usecases.GetCountryDataInput{Name: country, TargetLanguage: targetLang}
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
