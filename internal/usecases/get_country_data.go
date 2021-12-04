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
		en, ok := feature.Properties["name_en"].(string)
		if !ok {
			return nil, errors.Errorf("(%s) invalid english name type: %T", name, feature.Properties["name_en"])
		}
		prop, err := domain.NewProperties(name, tl, en, rubyLookup(tl))
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
	default:
		return "", errors.Errorf("unsupported language %q", in)
	}
}

func countryMapData(in string) (float64, float64, int) {
	switch in {
	case "japan":
		return 38.0, 137.0, 5
	case "united states of america":
		return 37.1, -95.7, 3
	default:
		return 0, 0, 1
	}
}

func rubyLookup(in string) string {
	switch in {
	case "鹿児島県":
		return "かごしまけん"
	case "大分県":
		return "おおいたけん"
	case "福岡県":
		return "ふくおかけん"
	case "佐賀県":
		return "さがけん"
	case "長崎県":
		return "ながさきけん"
	case "熊本県":
		return "くまもとけん"
	case "宮崎県":
		return "みやざきけん"
	case "徳島県":
		return "とくしまけん"
	case "香川県":
		return "かがわけん"
	case "愛媛県":
		return "えひめけん"
	case "高知県":
		return "こうちけん"
	case "島根県":
		return "しまねけん"
	case "山口県":
		return "やまぐちけん"
	case "鳥取県":
		return "とっとりけん"
	case "兵庫県":
		return "ひょうごけん"
	case "京都府":
		return "きょうとふ"
	case "福井県":
		return "ふくいけん"
	case "石川県":
		return "いしかわけん"
	case "富山県":
		return "とやまけん"
	case "新潟県":
		return "にいがたけん"
	case "山形県":
		return "やまがたけん"
	case "秋田県":
		return "あきたけん"
	case "青森県":
		return "あおもりけん"
	case "岩手県":
		return "いわてけん"
	case "宮城県":
		return "みやぎけん"
	case "福島県":
		return "ふくしまけん"
	case "茨城県":
		return "いばらきけん"
	case "千葉県":
		return "ちばけん"
	case "東京都":
		return "とうきょうと"
	case "神奈川県":
		return "かながわけん"
	case "静岡県":
		return "しずおかけん"
	case "愛知県":
		return "あいちけん"
	case "三重県":
		return "みえけん"
	case "和歌山県":
		return "わかやまけん"
	case "大阪府":
		return "おおさかふ"
	case "岡山県":
		return "おかやまけん"
	case "広島県":
		return "ひろしまけん"
	case "北海道":
		return "ほっかいどう"
	case "沖縄県":
		return "おきなわけん"
	case "群馬県":
		return "ぐんまけん"
	case "長野県":
		return "ながのけん"
	case "栃木県":
		return "とちぎけん"
	case "岐阜県":
		return "ぎふけん"
	case "滋賀県":
		return "しがけん"
	case "埼玉県":
		return "さいたまけん"
	case "山梨県":
		return "やまなしけん"
	case "奈良県":
		return "ならけん"
	default:
		return ""
	}
}
