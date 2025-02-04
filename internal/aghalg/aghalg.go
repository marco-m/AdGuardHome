// Package aghalg contains common generic algorithms and data structures.
//
// TODO(a.garipov): Move parts of this into golibs.
package aghalg

import (
	"fmt"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Coalesce returns the first non-zero value.  It is named after the function
// COALESCE in SQL.  If values or all its elements are empty, it returns a zero
// value.
func Coalesce[T comparable](values ...T) (res T) {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}

	return zero
}

// UniqChecker allows validating uniqueness of comparable items.
//
// TODO(a.garipov): The Ordered constraint is only really necessary in Validate.
// Consider ways of making this constraint comparable instead.
type UniqChecker[T constraints.Ordered] map[T]int64

// Add adds a value to the validator.  v must not be nil.
func (uc UniqChecker[T]) Add(elems ...T) {
	for _, e := range elems {
		uc[e]++
	}
}

// Merge returns a checker containing data from both uc and other.
func (uc UniqChecker[T]) Merge(other UniqChecker[T]) (merged UniqChecker[T]) {
	merged = make(UniqChecker[T], len(uc)+len(other))
	for elem, num := range uc {
		merged[elem] += num
	}

	for elem, num := range other {
		merged[elem] += num
	}

	return merged
}

// Validate returns an error enumerating all elements that aren't unique.
func (uc UniqChecker[T]) Validate() (err error) {
	var dup []T
	for elem, num := range uc {
		if num > 1 {
			dup = append(dup, elem)
		}
	}

	if len(dup) == 0 {
		return nil
	}

	slices.Sort(dup)

	return fmt.Errorf("duplicated values: %v", dup)
}
