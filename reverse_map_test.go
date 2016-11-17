package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func genMaps(mapSize, strLen int) (map[int]string, map[string]int, []string) {
	m := make(map[int]string)
	r := make(map[string]int)
	allstrings := []string{}
	for i := 0; i != mapSize; i++ {
		k, v := rand.Int(), RandStringBytesMaskImprSrc(strLen)
		m[k] = v
		r[v] = k
		allstrings = append(allstrings, v)
	}
	return m, r, allstrings
}

func BenchmarkLookup(b *testing.B) {
	for mapSize := 1; mapSize <= 2048; mapSize *= 8 {
		for strLen := 1; strLen < 30; strLen += 6 {
			testMap, rMap, allstrings := genMaps(mapSize, strLen)
			b.Run(
				fmt.Sprintf("mapSize=%d,strLen=%d", mapSize, strLen),
				func(b *testing.B) {
					b.Run("Iteration", func(b *testing.B) {
						for i := 0; i != b.N; i++ {
							for _, value := range allstrings {
								for _, v := range testMap {
									if v == value {
										break
									}
								}
							}
						}
					})
					b.Run("RMap", func(b *testing.B) {
						for i := 0; i != b.N; i++ {
							for _, value := range allstrings {
								_ = rMap[value]
							}
						}
					})
				})
		}
	}
}
