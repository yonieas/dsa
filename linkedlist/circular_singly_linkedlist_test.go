package linkedlist_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/linkedlist"
)

func TestCircularSinglyLinkedList(t *testing.T) {
	specs := []prop.Spec{
		prop.Append(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Prepend(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.GetSet(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.HeadTail(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Pop(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Shift(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TryPop(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TryShift(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TryHead(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TryTail(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TryGet(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TrySet(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TryRemove(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.TryCycle(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Iter(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Insert(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Remove(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Rotate(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.Cycle(linkedlist.NewCircularSinglyLinkedList[int]),
		prop.CircularIter(linkedlist.NewCircularSinglyLinkedList[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
