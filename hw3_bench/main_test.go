package main

import (
	"bytes"
	"io/ioutil"
	"testing"
	bnc "mail_course/hw3_bench/bench"
)

// запускаем перед основными функциями по разу чтобы файл остался в памяти в файловом кеше
// ioutil.Discard - это ioutil.Writer который никуда не пишет
func init() {
	bnc.SlowSearch(ioutil.Discard)
	bnc.FastSearch(ioutil.Discard)
}

// -----
// go test -v

func TestSearch(t *testing.T) {
	slowOut := new(bytes.Buffer)
	bnc.SlowSearch(slowOut)
	slowResult := slowOut.String()

	fastOut := new(bytes.Buffer)
	bnc.FastSearch(fastOut)
	fastResult := fastOut.String()

	if slowResult != fastResult {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", fastResult, slowResult)
	}
}

// -----
// go test -bench . -benchmem

func BenchmarkSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bnc.SlowSearch(ioutil.Discard)
	}
}

func BenchmarkFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bnc.FastSearch(ioutil.Discard)
	}
}
