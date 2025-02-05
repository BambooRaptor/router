package set

import (
	"fmt"
	"strings"
)

type Set[T comparable] map[T]interface{}

func FromArray[T comparable](arr []T) *Set[T] {
	s := New[T]()
	for _, el := range arr {
		s.Add(el)
	}
	return s
}

func New[T comparable]() *Set[T] {
	return &Set[T]{}
}

func (s *Set[T]) Has(el T) bool {
	_, ok := (*s)[el]
	return ok
}

func (s *Set[T]) Add(el T) {
	if !s.Has(el) {
		(*s)[el] = nil
	}
}

func (s *Set[T]) Remove(el T) { delete(*s, el) }

func (s *Set[T]) ToArray() []T {
	arr := make([]T, 0)
	for t := range *s {
		arr = append(arr, t)
	}
	return arr
}

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

func (s *Set[T]) String() string {
	arr := s.ToArray()
	els := make([]string, 0)
	for _, el := range arr {
		els = append(els, fmt.Sprintf("%v", el))
	}
	return strings.Join(els, ", ")
}
