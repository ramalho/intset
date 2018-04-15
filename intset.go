// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Extended by Luciano Ramalho according to exercises in the _GOPL_
// book. See `EXERCISES.adoc`.

// Package intset provides a set of integers based on a bit vector.
package intset

import (
	"bytes"
	"strconv"
	"strings"
)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
	len   int
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && (s.words[word]>>bit)&1 == 1
}

// Add the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
	s.len++
}

// Len reports the number of elements in the set.
// Also known as set cardinality.
func (s *IntSet) Len() int {
	return s.len
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if tword > 0 {
			if i < len(s.words) {
				before := bitCount(s.words[i])
				s.words[i] |= tword
				s.len += bitCount(s.words[i]) - before
			} else {
				s.words = append(s.words, tword)
				s.len += bitCount(tword)
			}
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	buf.WriteString(strings.Join(s.elemStr(), " "))
	buf.WriteByte('}')
	return buf.String()
}

// NewFromSlice returns pointer to a new IntSet with the
// elements given in the slice
func NewFromSlice(slice []int) *IntSet {
	s := IntSet{}
	for _, n := range slice {
		s.Add(n)
	}
	return &s
}

func bitCount(word uint64) int {
	count := 0
	for bit := uint(0); bit < 64; bit++ {
		count += int(word>>bit) & 1
	}
	return count
}

// Elems returns a new slice of integers with the elements
// found in the IntSet in ascending order.
func (s *IntSet) Elems() []int {
	elems := []int{}
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				elems = append(elems, 64*i+j)
			}
		}
	}
	return elems
}

func (s *IntSet) elemStr() []string {
	elems := s.Elems()
	res := make([]string, len(elems))
	for i, v := range elems {
		res[i] = strconv.Itoa(v)
	}
	return res
}
