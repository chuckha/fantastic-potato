package web

import (
	"encoding/json"
	"net/http"

	"github.com/chuckha/geogame6/internal/usecases"
	"github.com/pkg/errors"
)

type presenter struct{}

type OutData struct {
	GeoJSON     interface{} `json:"geojson"`
	CenterLat   float64     `json:"center_lat"`
	CenterLon   float64     `json:"center_lon"`
	DefaultZoom int         `json:"default_zoom"`
}

func (p *presenter) presentGetCountryData(w http.ResponseWriter, out *usecases.GetCountryDataOutput) error {
	w.Header().Set("Content-Type", "application/json")
	o := &OutData{
		GeoJSON:     string(out.RawData),
		CenterLat:   out.CenterLat,
		CenterLon:   out.CenterLon,
		DefaultZoom: out.DefaultZoom,
	}
	data, err := json.Marshal(o)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := w.Write(data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
