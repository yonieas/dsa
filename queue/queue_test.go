package queue_test

import (
	"math/rand"
	"testing"

	"github.com/josestg/dsa/adt/adttest"
	"github.com/josestg/dsa/queue"
)

func TestQueue(t *testing.T) {
	c := queue.New[int]
	g := func() int {
		return rand.Intn(128)
	}
	simulator := adttest.QueueSimulator(c, g)
	simulator.Run(t)
}
