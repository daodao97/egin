package _switch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOn(t *testing.T) {
	s := NewSwitch("http://192.168.111.6:9500")
	assert.True(t, s.IsOn("oms/coupon_maintenance"))
}

// BenchmarkIsOn-8   	      36	  28837906 ns/op	   67312 B/op	     168 allocs/op
func BenchmarkIsOn(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewSwitch("http://192.168.111.6:9500").IsOn("oms/coupon_maintenance")
	}
}
