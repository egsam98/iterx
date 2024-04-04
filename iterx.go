//go:build goexperiment.rangefunc

package iterx

import (
	"encoding/json"
	"io"
	"iter"

	"github.com/egsam98/errors"
	"gopkg.in/yaml.v3"
)

type Result[T any] struct {
	Ok  T
	Err error
}

type Validator interface {
	Validate() error
}

func JSON[T Validator](r io.Reader) (iter.Seq2[int, Result[T]], error) {
	var arr []json.RawMessage
	if err := json.NewDecoder(r).Decode(&arr); err != nil {
		return nil, errors.Wrapf(err, "JSON array is expected")
	}

	return func(yield func(int, Result[T]) bool) {
		for i, elem := range arr {
			var t T
			var err error
			if err = json.Unmarshal(elem, &t); err == nil {
				err = t.Validate()
			}
			if !yield(i, Result[T]{Ok: t, Err: err}) {
				return
			}
		}
	}, nil
}

func YAML[T Validator](r io.Reader) (iter.Seq2[int, Result[T]], error) {
	var arr []yaml.Node
	if err := yaml.NewDecoder(r).Decode(&arr); err != nil {
		return nil, errors.Wrapf(err, "JSON array is expected")
	}

	return func(yield func(int, Result[T]) bool) {
		for i, node := range arr {
			var t T
			var err error
			if err = node.Decode(&t); err == nil {
				err = t.Validate()
			}
			if !yield(i, Result[T]{Ok: t, Err: err}) {
				return
			}
		}
	}, nil
}
