package linkedlist_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/linkedlist"
)

func TestSinglyLinkedList(t *testing.T) {
	specs := []prop.Spec{
		prop.Append(linkedlist.NewSinglyLinkedList[int]),
		prop.Prepend(linkedlist.NewSinglyLinkedList[int]),
		prop.GetSet(linkedlist.NewSinglyLinkedList[int]),
		prop.HeadTail(linkedlist.NewSinglyLinkedList[int]),
		prop.Pop(linkedlist.NewSinglyLinkedList[int]),
		prop.Shift(linkedlist.NewSinglyLinkedList[int]),
		prop.TryPop(linkedlist.NewSinglyLinkedList[int]),
		prop.TryShift(linkedlist.NewSinglyLinkedList[int]),
		prop.TryHead(linkedlist.NewSinglyLinkedList[int]),
		prop.TryTail(linkedlist.NewSinglyLinkedList[int]),
		prop.TryGet(linkedlist.NewSinglyLinkedList[int]),
		prop.TrySet(linkedlist.NewSinglyLinkedList[int]),
		prop.TryRemove(linkedlist.NewSinglyLinkedList[int]),
		prop.Iter(linkedlist.NewSinglyLinkedList[int]),
		prop.IterBackward(linkedlist.NewSinglyLinkedList[int]),
		prop.Insert(linkedlist.NewSinglyLinkedList[int]),
		prop.Remove(linkedlist.NewSinglyLinkedList[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
