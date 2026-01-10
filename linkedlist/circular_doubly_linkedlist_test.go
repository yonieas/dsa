package linkedlist_test

import (
	"slices"
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/linkedlist"
)

func TestCircularDoublyLinkedList(t *testing.T) {
	specs := []prop.Spec{
		prop.Append(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Prepend(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.GetSet(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.HeadTail(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Pop(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Shift(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryPop(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryShift(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryHead(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryTail(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryGet(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TrySet(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryRemove(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryCycle(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.TryReverseCycle(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Iter(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.IterBackward(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Insert(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Remove(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Rotate(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.Cycle(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.ReverseCycle(linkedlist.NewCircularDoublyLinkedList[int]),
		prop.CircularIter(linkedlist.NewCircularDoublyLinkedList[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}

func TestCircularDoublyLinkedList_CircularIteratorBackward(t *testing.T) {
	l := linkedlist.NewCircularDoublyLinkedList[int]()
	l.Append(1)
	l.Append(2)
	l.Append(3)

	var collected []int
	count := 0
	for v := range l.CircularIteratorBackward {
		collected = append(collected, v)
		count++
		if count >= 7 {
			break
		}
	}
	want := []int{3, 2, 1, 3, 2, 1, 3}
	if !slices.Equal(collected, want) {
		t.Errorf("got %v, want %v", collected, want)
	}
}
