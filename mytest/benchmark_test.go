package main

import (
	"testing"
)

func BenchmarkClusterLargeValue(b *testing.B) {
	if !request(largeValLen, serverSetUrls[0:serverCnt], b) {
		b.Error("cluster benchmark with small value length failed")
	}
}
func BenchmarkClusterLargeValueAsync(b *testing.B) {

	if !asyncRequest(largeValLen, serverSetUrls[0:serverCnt], b) {
		b.Error("async cluster benchmark with small value length failed")
	}
}

func BenchmarkLargeValue(b *testing.B) {
	if !request(largeValLen, serverSetUrls[0:1], b) {
		b.Error("benchmark for signal node with large value length failed")
	}
}

func BenchmarkLargeValueAsync(b *testing.B) {
	if !asyncRequest(largeValLen, serverSetUrls[0:1], b) {
		b.Error("async benchmark for signal node with large value length failed")
	}
}

func BenchmarkClusterSmallValue(b *testing.B) {
	if !request(smallValLen, serverSetUrls[0:serverCnt], b) {
		b.Error("cluster benchmark with small value length failed")
	}
}

func BenchmarkClusterSmallValueAsync(b *testing.B) {

	if !asyncRequest(smallValLen, serverSetUrls[0:serverCnt], b) {
		b.Error("async cluster benchmark with small value length failed")
	}
}

func BenchmarkSmallValue(b *testing.B) {
	if !request(smallValLen, serverSetUrls[0:1], b) {
		b.Error("benchmark for signal node with small value length failed")
	}
}
func BenchmarkSmallValueAsync(b *testing.B) {

	if !asyncRequest(smallValLen, serverSetUrls[0:1], b) {
		b.Error("async benchmark for signal node with small value length failed")
	}
}
