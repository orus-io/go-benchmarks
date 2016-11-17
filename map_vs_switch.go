package main

import (
	"fmt"
	"math/rand"
	"os"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type typeValues map[int]string

var types = []struct {
	name   string
	values typeValues
}{
	{
		"one value",
		typeValues{0: "Monday"},
	},
	{
		"days (7)",
		typeValues{
			0: "Monday",
			1: "Tuesday",
			2: "Wednesday",
			3: "Thursday",
			4: "Friday",
			5: "Satursday",
			6: "Sunday",
		},
	},
	{
		"prime (13)",
		typeValues{
			2:  "p2",
			3:  "p3",
			5:  "p5",
			7:  "p7",
			11: "p11",
			13: "p13",
			17: "p17",
			19: "p19",
			23: "p23",
			29: "p29",
			31: "p31",
			41: "p41",
			43: "p43",
		},
	},
	{
		"26 values",
		typeValues{
			2:   "p2",
			3:   "p3",
			5:   "p5",
			7:   "p7",
			11:  "p11",
			13:  "p13",
			17:  "p17",
			19:  "p19",
			23:  "p23",
			29:  "p29",
			31:  "p31",
			41:  "p41",
			43:  "p43",
			52:  "p52",
			53:  "p53",
			55:  "p55",
			57:  "p57",
			511: "p511",
			513: "p513",
			517: "p517",
			519: "p519",
			523: "p5523",
			529: "p529",
			531: "p531",
			541: "p541",
			543: "p543",
		},
	},
}

func main() {
	f, err := os.Create("map_vs_switch_test.go")
	check(err)
	defer f.Close()

	w := func(format string, args ...interface{}) {
		fmt.Fprintf(f, format, args...)
	}
	w(header)

	for i := 32; i <= 1024; i *= 2 {
		values := make(typeValues)
		for k := 0; k < i; k++ {
			values[rand.Int()] = RandStringBytesMaskImprSrc(rand.Intn(30) + 2)
		}
		types = append(types, struct {
			name   string
			values typeValues
		}{
			fmt.Sprintf("random %d", i),
			values,
		})
	}

	for _, t := range types {
		w("\tb.Run(\"%s\", func(b *testing.B) {\n", t.name)

		w("\t\tm := map[string]int{\n")
		for k, v := range t.values {
			w("\t\t\t\"%s\": %d,\n", v, k)
		}
		w("\t\t}\n")

		w("\t\tmapLookup := func(value string) (int, error) {\n")
		w("\t\t\tv, ok := m[value]\n")
		w("\t\t\tif !ok {\n")
		w("\t\t\t\treturn 0, errors.New(\"invalid\")\n")
		w("\t\t\t}\n")
		w("\t\t\treturn v, nil\n")
		w("\t\t}\n")

		w("\t\tswitchLookup := func(value string) (int, error) {\n")
		w("\t\t\tswitch(value) {\n")
		for k, v := range t.values {
			w("\t\t\tcase \"%s\":\n", v)
			w("\t\t\t\treturn %d, nil\n", k)
		}
		w("\t\t\tdefault:\n")
		w("\t\t\t\treturn 0, errors.New(\"invalid\")\n")
		w("\t\t\t}\n")
		w("\t\t}\n")

		w("\t\tb.Run(\"map\", func(b *testing.B) {\n")
		w("\t\t\tfor i := 0; i != b.N; i++ {\n")
		for _, v := range t.values {
			w("\t\t\t\t_, _ = mapLookup(\"%s\")\n", v)
		}
		w("\t\t\t\t\n")
		w("\t\t\t}\n")
		w("\t})\n")

		w("\t\tb.Run(\"switch\", func(b *testing.B) {\n")
		w("\t\t\tfor i := 0; i != b.N; i++ {\n")
		for _, v := range t.values {
			w("\t\t\t\t_, _ = switchLookup(\"%s\")\n", v)
		}
		w("\t\t\t\t\n")
		w("\t\t\t}\n")

		w("\t\t})\n")
		w("\t})\n")
	}
	w("}\n")
}

const header = `package main

import (
	"errors"
	"testing"
)

func BenchmarkMapVsSwitch(b *testing.B) {
`

/*
	b.Run(
		"3 element",
		func(b *testing.B) {
			m := make(map[string]int)
			m["something"] = 42
			m["other"] = 43
			m["plouf"] = 44
			values := []string{"something", "other", "plouf"}
			b.Run("map", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for _, v := range values {
						_, _ = m[v]
					}
				}
			})
			b.Run("switch", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for _, v := range values {
						switch v {
						case "something":
							_ = 42
						case "other":
							_ = 43
						case "plouf":
							_ = 44
						default:
							_ = false
						}
					}
				}
			})
		},
	)
*/
