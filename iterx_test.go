//go:build goexperiment.rangefunc

package iterx

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	Label string `json:"label"`
}

func (a *A) Validate() error {
	return nil
}

func TestJSON(t *testing.T) {
	seq, err := JSON[*A](jsonReader())
	assert.NoError(t, err)

	for i, res := range seq {
		assert.NoError(t, res.Err)
		t.Log(i, res.Ok)
	}
}

type Data struct {
	Identifier string `yaml:"identifier"`
	Title      string `yaml:"title"`
}

func (d *Data) Validate() error {
	return nil
}

func TestYAML(t *testing.T) {
	seq, err := YAML[*Data](yamlReader())
	assert.NoError(t, err)

	for i, res := range seq {
		assert.NoError(t, res.Err)
		t.Log(i, res.Ok)
	}
}

func jsonReader() io.Reader {
	f, err := os.Open("test.json")
	if err != nil {
		panic(err)
	}
	return f
}

func yamlReader() io.Reader {
	f, err := os.Open("test.yaml")
	if err != nil {
		panic(err)
	}
	return f
}
