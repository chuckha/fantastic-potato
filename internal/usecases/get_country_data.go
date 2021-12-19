package usecases

import (
	"context"
	"strings"

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
	Name           string
	TargetLanguage string
}
type GetCountryDataOutput struct {
	CenterLat   float64
	CenterLon   float64
	DefaultZoom int
	RawData     []byte
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
	languagePropertyKey, err := languageToProperty(in.TargetLanguage)
	if err != nil {
		return nil, err
	}

	for i, feature := range fc1.Features {
		name, ok := feature.Properties["name"].(string)
		if !ok {
			return nil, errors.Errorf("invalid feature name type: %T", feature.Properties["name"])
		}
		tl, ok := feature.Properties[languagePropertyKey].(string)
		if !ok {
			return nil, errors.Errorf("(%s) invalid japanese name type: %T", name, feature.Properties["name_jp"])
		}
		tl = JPClean(tl)
		prop, err := domain.NewProperties(name, tl, rubyLookup(tl))
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
	lat, lon, zoom := countryMapData(in.Name)
	return &GetCountryDataOutput{
		RawData:     cleanData,
		CenterLat:   lat,
		CenterLon:   lon,
		DefaultZoom: zoom,
	}, nil
}

func languageToProperty(in string) (string, error) {
	switch in {
	case "japanese":
		return "name_ja", nil
	case "english":
		return "name_en", nil
	case "french":
		return "name_fr", nil
	case "korean":
		return "name_ko", nil
	default:
		return "", errors.Errorf("unsupported language %q", in)
	}
}

func countryMapData(in string) (float64, float64, int) {
	switch in {
	case "japan":
		return 38.0, 137.0, 5
	case "korea":
		return 36.49869367154783, 127.91915640742643, 6
	case "united states of america":
		return 37.1, -95.7, 3
	default:
		return 0, 0, 1
	}
}

func rubyLookup(in string) string {
	switch in {
	case "鹿児島":
		return "かごしま"
	case "大分":
		return "おおいた"
	case "福岡":
		return "ふくおか"
	case "佐賀":
		return "さが"
	case "長崎":
		return "ながさき"
	case "熊本":
		return "くまもと"
	case "宮崎":
		return "みやざき"
	case "徳島":
		return "とくしま"
	case "香川":
		return "かがわ"
	case "愛媛":
		return "えひめ"
	case "高知":
		return "こうち"
	case "島根":
		return "しまね"
	case "山口":
		return "やまぐち"
	case "鳥取":
		return "とっとり"
	case "兵庫":
		return "ひょうご"
	case "京都":
		return "きょうと"
	case "福井":
		return "ふくい"
	case "石川":
		return "いしかわ"
	case "富山":
		return "とやま"
	case "新潟":
		return "にいがた"
	case "山形":
		return "やまがた"
	case "秋田":
		return "あきた"
	case "青森":
		return "あおもり"
	case "岩手":
		return "いわて"
	case "宮城":
		return "みやぎ"
	case "福島":
		return "ふくしま"
	case "茨城":
		return "いばらき"
	case "千葉":
		return "ちば"
	case "東京":
		return "とうきょう"
	case "神奈川":
		return "かながわ"
	case "静岡":
		return "しずおか"
	case "愛知":
		return "あいち"
	case "三重":
		return "みえ"
	case "和歌山":
		return "わかやま"
	case "大阪":
		return "おおさか"
	case "岡山":
		return "おかやま"
	case "広島":
		return "ひろしま"
	case "北海道":
		return "ほっかいどう"
	case "沖縄":
		return "おきなわ"
	case "群馬":
		return "ぐんま"
	case "長野":
		return "ながの"
	case "栃木":
		return "とちぎ"
	case "岐阜":
		return "ぎふ"
	case "滋賀":
		return "しが"
	case "埼玉":
		return "さいたま"
	case "山梨":
		return "やまなし"
	case "奈良":
		return "なら"
	default:
		return ""
	}
}

func JPClean(s string) string {
	// special cases
	if s == "京都" {
		return s
	}
	return strings.Trim(s, "県府都")
}
