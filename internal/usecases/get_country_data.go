package usecases

import (
	"context"

	"github.com/chuckha/geogame6/internal/domain"
	geojson "github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
)

type gatherer interface {
	Gather(ctx context.Context, name string) ([]byte, error)
}

type encdec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type CountryData struct {
	Gatherer gatherer
	EncDec   encdec
}

func NewCountryData(gatherer gatherer, encdec encdec) *CountryData {
	return &CountryData{
		Gatherer: gatherer,
		EncDec:   encdec,
	}
}

type GetCountryDataInput struct {
	Name string
}
type GetCountryDataOutput struct {
	RawData []byte
}

func (g *CountryData) GetCountryData(ctx context.Context, in *GetCountryDataInput) (*GetCountryDataOutput, error) {
	data, err := g.Gatherer.Gather(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	// get the data into actual geojson then build up quiz data obj
	fc1, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	features := make([]*domain.Feature, len(fc1.Features))
	for i, feature := range fc1.Features {
		name, ok := feature.Properties["name"].(string)
		if !ok {
			return nil, errors.Errorf("invalid feature name type: %T", feature.Properties["name"])
		}
		jp, ok := feature.Properties["name_ja"].(string)
		if !ok {
			return nil, errors.Errorf("(%s) invalid japanese name type: %T", name, feature.Properties["name_jp"])
		}
		en, ok := feature.Properties["name_en"].(string)
		if !ok {
			return nil, errors.Errorf("(%s) invalid english name type: %T", name, feature.Properties["name_en"])
		}
		prop, err := domain.NewProperties(name, jp, en, "")
		if err != nil {
			return nil, err
		}
		df, err := domain.NewFeature(feature.Type, prop, feature.Geometry)
		if err != nil {
			return nil, err
		}
		features[i] = df
	}
	qd, err := domain.NewQuizData(fc1.Type, in.Name, features)
	if err != nil {
		return nil, err
	}
	cleanData, err := g.EncDec.Encode(qd)
	if err != nil {
		return nil, err
	}
	return &GetCountryDataOutput{
		RawData: cleanData,
	}, nil
}
