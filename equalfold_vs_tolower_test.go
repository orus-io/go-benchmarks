package main

import (
	"strings"
	"testing"
)

var (
	s1 = "A sTrInG"
	s2 = "a StRiNg"
)

func BenchmarkLower(b *testing.B) {
	ls1 := strings.ToLower(s1)
	for i := 0; i != b.N; i++ {
		_ = ls1 == strings.ToLower(s2)
	}
}

func BenchmarkEqualFold(b *testing.B) {
	for i := 0; i != b.N; i++ {
		_ = strings.EqualFold(s1, s2)
	}
}
