package web

import (
	"net/http"

	"github.com/chuckha/geogame6/internal/usecases"
	"github.com/pkg/errors"
)

type presenter struct{}

func (p *presenter) presentGetCountryData(w http.ResponseWriter, out *usecases.GetCountryDataOutput) error {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(out.RawData); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
