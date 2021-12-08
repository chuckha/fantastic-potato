package domain

import (
	"errors"

	geojson "github.com/paulmach/go.geojson"
)

// QuizData must be valid geojson.
// However properties are severely limited from traditional GeoJSON.
type QuizData struct {
	Type     string     `json:"type"`
	Name     string     `json:"name"`
	Features []*Feature `json:"features"`
}

func NewQuizData(geotype, name string, features []*Feature) (*QuizData, error) {
	if geotype == "" {
		return nil, errors.New("quiz data requires a type")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}
	if len(features) == 0 {
		return nil, errors.New("a quiz must have at least one feature")
	}
	return &QuizData{
		Type:     geotype,
		Name:     name,
		Features: features,
	}, nil
}

type Feature struct {
	Type       string            `json:"type"`
	Properties *Properties       `json:"properties"`
	Geometry   *geojson.Geometry `json:"geometry"`
}

func NewFeature(featureType string, prop *Properties, geo *geojson.Geometry) (*Feature, error) {
	if featureType == "" {
		return nil, errors.New("a feature requires a non empty type")
	}
	return &Feature{
		Type:       featureType,
		Properties: prop,
		Geometry:   geo,
	}, nil
}

type Properties struct {
	ID                 string `json:"name"`
	TargetLanguageName string `json:"name_tl"`
	RubyText           string `json:"ruby_text"` // optional
}

func NewProperties(id, tl, ruby string) (*Properties, error) {
	if id == "" {
		return nil, errors.New("properties requires an id")
	}
	if tl == "" {
		return nil, errors.New("properties must have a tl name")
	}
	return &Properties{
		ID:                 id,
		TargetLanguageName: tl,
		RubyText:           ruby,
	}, nil
}
