package stringutil

import "testing"

func defaultBytesToString(bytes []byte) string {
	return string(bytes)
}

func defaultStringToBytes(s string) []byte {
	return []byte(s)
}

func BenchmarkBytesToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BytesToString(bytes)
	}
}

func BenchmarkDefaultBytesToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defaultBytesToString(bytes)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringToBytes(str)
	}
}

func BenchmarkDefaultStringToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defaultStringToBytes(str)
	}
}
