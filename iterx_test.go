//go:build goexperiment.rangefunc

package iterx

import (
	"io"
	"iter"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceErr(t *testing.T) {
	var seq iter.Seq2[int, error] = func(yield func(int, error) bool) {
		for i := range 3 {
			yield(i, nil)
		}
	}
	s, err := SliceErr(seq)
	assert.NoError(t, err)
	t.Log(s)
}

func TestMapErr(t *testing.T) {
	var seq iter.Seq2[int, error] = func(yield func(int, error) bool) {
		for i := range 3 {
			yield(i, nil)
		}
	}
	s, err := MapErr(seq, func(v int) string {
		return strconv.Itoa(v)
	})
	assert.NoError(t, err)
	t.Log(s)
}

type A struct {
	Label string `json:"label"`
}

func (a *A) Validate() error {
	return nil
}

func TestJSON(t *testing.T) {
	seq, err := JSON[*A](jsonReader())
	assert.NoError(t, err)

	for res := range seq {
		assert.NoError(t, res.Err)
		t.Log(res.Index, res.Data)
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

	for res := range seq {
		assert.NoError(t, res.Err)
		t.Log(res.Index, res.Data)
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
