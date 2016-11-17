package main

import (
	"bytes"
	"testing"
)

var (
	s1 = "a string"
	s2 = "another string"
	b1 = []byte("a string")
	b2 = []byte("another string")
)

func BenchmarkByteToStringCmp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = s1 == string(b1)
		_ = s2 == string(b1)
		_ = s1 == string(b2)
		_ = s2 == string(b2)
	}
}

func BenchmarkByteToStringOnceCmp(b *testing.B) {
	sb1 := string(b1)
	sb2 := string(b2)
	for i := 0; i < b.N; i++ {
		_ = s1 == sb1
		_ = s2 == sb1
		_ = s1 == sb2
		_ = s2 == sb2
	}
}

func BenchmarkStringToByteCmp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes.Equal([]byte(s1), b1)
		bytes.Equal([]byte(s2), b1)
		bytes.Equal([]byte(s1), b2)
		bytes.Equal([]byte(s2), b2)
	}
}

func BenchmarkStringToByteOnceCmp(b *testing.B) {
	bs1 := []byte(s1)
	bs2 := []byte(s2)
	for i := 0; i < b.N; i++ {
		bytes.Equal(bs1, b1)
		bytes.Equal(bs1, b1)
		bytes.Equal(bs2, b2)
		bytes.Equal(bs2, b2)
	}
}
