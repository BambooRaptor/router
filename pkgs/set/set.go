package set

import (
	"fmt"
	"strings"
)

type Set[T comparable] map[T]interface{}

// Create a set from an existing array
func FromArray[T comparable](arr []T) Set[T] {
	s := New[T]()
	for _, el := range arr {
		_ = s.Add(el)
	}
	return s
}

// Create a new, empty set
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Check if the set has a particular element
func (s *Set[T]) Has(el T) bool {
	_, ok := (*s)[el]
	return ok
}

// Add an element to the set.
// Errors out if it doesn't
func (s *Set[T]) Add(el T) error {
	if s.Has(el) {
		return fmt.Errorf("set already has element %v", el)
	}
	(*s)[el] = nil
	return nil
}

// Removes an element from the set
func (s *Set[T]) Remove(el T) { delete(*s, el) }

// Compile the set into an array
func (s *Set[T]) ToArray() []T {
	arr := make([]T, 0)
	for t := range *s {
		arr = append(arr, t)
	}
	return arr
}

// Compares two sets and checks if all the elements match
func (s *Set[T]) Matches(ns *Set[T]) bool {
	ls := s.ToArray()
	lns := ns.ToArray()

	if len(ls) != len(lns) {
		return false
	}

	for _, el := range ls {
		if !ns.Has(el) {
			return false
		}
	}

	return true
}

// For debugging purposes.
// If you want to create your own string method,
// you can use the ToArray() function.
// Or create a type alias
func (s *Set[T]) String() string {
	arr := s.ToArray()
	els := make([]string, 0)
	for _, el := range arr {
		els = append(els, fmt.Sprintf("%v", el))
	}
	return strings.Join(els, ", ")
}
