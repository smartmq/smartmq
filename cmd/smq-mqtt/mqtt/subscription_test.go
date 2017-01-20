package mqtt

import (
	"testing"
)

func BenchmarkSubscription_IsSubscribedEQ(b *testing.B) {
	s := NewSubscription("/test/j/0000/c/1021", 0x00)
	for n := 0; n < b.N; n++ {
		s.IsSubscribed("/test/j/0000/c/1021")
	}
}

func BenchmarkSubscription_IsSubscribedWC(b *testing.B) {
	s := NewSubscription("/+/j/+/c/#", 0x00)
	for n := 0; n < b.N; n++ {
		s.IsSubscribed("/test/j/0000/c/1021")
	}
}

func TestSubscription_IsSubscribed(t *testing.T) {
	s1 := NewSubscription("/+/j/+/c/#", 0x00)
	b1 := s1.IsSubscribed("/test/j/0000/c/1021")
	if !b1 {
		t.Fail()
	}

	s2 := NewSubscription("a", 0x00)
	b2 := s2.IsSubscribed("a")
	if !b2 {
		t.Fail()
	}
}
