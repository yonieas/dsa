package linkedlist_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/linkedlist"
)

func TestDoublyLinkedList(t *testing.T) {
	specs := []prop.Spec{
		prop.Append(linkedlist.NewDoublyLinkedList[int]),
		prop.Prepend(linkedlist.NewDoublyLinkedList[int]),
		prop.GetSet(linkedlist.NewDoublyLinkedList[int]),
		prop.HeadTail(linkedlist.NewDoublyLinkedList[int]),
		prop.Pop(linkedlist.NewDoublyLinkedList[int]),
		prop.Shift(linkedlist.NewDoublyLinkedList[int]),
		prop.TryPop(linkedlist.NewDoublyLinkedList[int]),
		prop.TryShift(linkedlist.NewDoublyLinkedList[int]),
		prop.TryHead(linkedlist.NewDoublyLinkedList[int]),
		prop.TryTail(linkedlist.NewDoublyLinkedList[int]),
		prop.TryGet(linkedlist.NewDoublyLinkedList[int]),
		prop.TrySet(linkedlist.NewDoublyLinkedList[int]),
		prop.TryRemove(linkedlist.NewDoublyLinkedList[int]),
		prop.Iter(linkedlist.NewDoublyLinkedList[int]),
		prop.IterBackward(linkedlist.NewDoublyLinkedList[int]),
		prop.Insert(linkedlist.NewDoublyLinkedList[int]),
		prop.Remove(linkedlist.NewDoublyLinkedList[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
