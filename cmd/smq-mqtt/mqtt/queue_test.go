package mqtt

import (
	"testing"
)

func TestQueue_EnqueueMessage(t *testing.T) {
	q := NewQueue()
	q.EnqueueMessage("uno")
	q.EnqueueMessage("due")
	q.EnqueueMessage("tre")

	if q.DequeueMessage() != "uno" {
		t.Fail()
	}
	if q.DequeueMessage() != "due" {
		t.Fail()
	}
	if q.DequeueMessage() != "tre" {
		t.Fail()
	}
}

func TestQueue_DequeueMessage(t *testing.T) {
	q := NewQueue()
	q.EnqueueMessage("uno")

	if q.DequeueMessage() != "uno" {
		t.Fail()
	}

	if q.DequeueMessage() != nil {
		t.Fail()
	}
}
