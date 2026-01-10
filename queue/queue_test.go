package queue_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/dynamicarray"
	"github.com/josestg/dsa/queue"
)

func TestQueue(t *testing.T) {
	specs := []prop.Spec{
		prop.Queue(queue.New[int]),
		prop.TryPeekQueue(queue.New[int]),
		prop.TryDequeue(queue.New[int]),
		prop.IterQueue(queue.New[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}

func TestQueueWithArrayBackend(t *testing.T) {
	f := func() *queue.Queue[int] {
		return queue.NewWith(dynamicarray.New[int](4))
	}
	specs := []prop.Spec{
		prop.Queue(f),
		prop.TryPeekQueue(f),
		prop.TryDequeue(f),
		prop.IterQueue(f),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
