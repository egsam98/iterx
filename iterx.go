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
	Index uint
	Data  T
	Err   error
}

type Validator interface {
	Validate() error
}

func JSON[T Validator](r io.Reader) (iter.Seq[Result[T]], error) {
	var arr []json.RawMessage
	if err := json.NewDecoder(r).Decode(&arr); err != nil {
		return nil, errors.Wrapf(err, "JSON array is expected")
	}

	return func(yield func(Result[T]) bool) {
		for i, elem := range arr {
			var t T
			var err error
			if err = json.Unmarshal(elem, &t); err == nil {
				err = t.Validate()
			}
			if !yield(Result[T]{
				Index: uint(i),
				Data:  t,
				Err:   err,
			}) {
				return
			}
		}
	}, nil
}

func YAML[T Validator](r io.Reader) (iter.Seq[Result[T]], error) {
	var arr []yaml.Node
	if err := yaml.NewDecoder(r).Decode(&arr); err != nil {
		return nil, errors.Wrapf(err, "JSON array is expected")
	}

	return func(yield func(Result[T]) bool) {
		for i, node := range arr {
			var t T
			var err error
			if err = node.Decode(&t); err == nil {
				err = t.Validate()
			}
			if !yield(Result[T]{
				Index: uint(i),
				Data:  t,
				Err:   err,
			}) {
				return
			}
		}
	}, nil
}

func SliceErr[T any](seq iter.Seq2[T, error]) ([]T, error) {
	var slice []T
	for t, err := range seq {
		if err != nil {
			return nil, err
		}
		slice = append(slice, t)
	}
	return slice, nil
}

func MapErr[K comparable, V any](seq iter.Seq2[V, error], keyFn func(V) K) (map[K]V, error) {
	m := make(map[K]V)
	for v, err := range seq {
		if err != nil {
			return nil, err
		}
		m[keyFn(v)] = v
	}
	return m, nil
}
