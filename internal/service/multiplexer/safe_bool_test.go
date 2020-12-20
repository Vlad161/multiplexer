package multiplexer_test

import (
	"testing"

	"github.com/vlad161/multiplexer/internal/service/multiplexer"
)

/*
	BenchmarkSafeBoolMutex-8        	10732446	       113 ns/op
	BenchmarkSafeBoolAtomic-8       	27321654	        40.2 ns/op
	BenchmarkSafeBoolAtomicType-8   	37224340	        32.5 ns/op
*/

func BenchmarkSafeBoolMutex(b *testing.B) {
	v := false
	sb := multiplexer.SafeBoolMutex{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sb.Get()
			sb.Set(v)
			v = !v
		}
	})
}

func BenchmarkSafeBoolAtomic(b *testing.B) {
	v := false
	sb := multiplexer.SafeBoolAtomic{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sb.Get()
			sb.Set(v)
			v = !v
		}
	})
}

func BenchmarkSafeBoolAtomicType(b *testing.B) {
	v := false
	sb := multiplexer.NewSafeBoolAtomicType()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sb.Get()
			sb.Set(v)
			v = !v
		}
	})
}
