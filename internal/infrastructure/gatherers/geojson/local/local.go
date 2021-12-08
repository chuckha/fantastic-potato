package local

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Gatherer struct {
	StatesFile string
	Outdir     string
}

func NewGatherer(statesFile, outdir string) *Gatherer {
	return &Gatherer{StatesFile: statesFile, Outdir: outdir}
}

/*
ogr2ogr -f GeoJSON states.json ne_10m_admin_1_states_provinces.shp
ogr2ogr -where "iso_a2='JP'" prefectures.json states.json
*/
func (g *Gatherer) Gather(ctx context.Context, name string) ([]byte, error) {
	outfile := g.outfile(name)
	if _, err := os.Stat(outfile); err == nil {
		out, err := os.ReadFile(outfile)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return out, nil
	}
	iso, err := countryToISO(name)
	if err != nil {
		return nil, err
	}
	os.MkdirAll(filepath.Dir(outfile), os.ModeDir|0o755)
	cmd := exec.CommandContext(ctx, "ogr2ogr", []string{
		"-where",
		fmt.Sprintf("iso_a2='%s'", iso),
		outfile,
		g.StatesFile}...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return nil, errors.WithStack(err)
	}
	data, err := ioutil.ReadFile(outfile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (g *Gatherer) outfile(country string) string {
	return filepath.Join(g.Outdir, country, "data.json")
}

func countryToISO(country string) (string, error) {
	switch strings.ToLower(country) {
	case "japan":
		return "JP", nil
	case "korea":
		return "KR", nil
	case "united states of america":
		return "US", nil
	default:
		return "", errors.Errorf("unknown country %q", country)
	}
}
