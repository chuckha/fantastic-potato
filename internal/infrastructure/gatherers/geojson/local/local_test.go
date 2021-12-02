package local

import (
	"context"
	"testing"
)

func TestLocalGather(t *testing.T) {
	gatherer := NewGatherer("../../../../../data/states.json", "./testdata")
	data, err := gatherer.Gather(context.Background(), "united states of america")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if len(data) == 0 {
		t.Fatal("data should not be empty")
	}
}
