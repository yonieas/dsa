package tree_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/tree"
)

func TestBinarySearchTree(t *testing.T) {
	specs := []prop.Spec{
		prop.AddExistsDel(tree.NewBinarySearchTree[int]),
		prop.BSTMinMax(tree.NewBinarySearchTree[int]),
		prop.BSTInOrder(tree.NewBinarySearchTree[int]),
		prop.BSTPreOrder(tree.NewBinarySearchTree[int]),
		prop.BSTPostOrder(tree.NewBinarySearchTree[int]),
		prop.BSTString(tree.NewBinarySearchTree[int]),
		prop.BSTIterBackward(tree.NewBinarySearchTree[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}

func TestBinarySearchTree_String_StringType(t *testing.T) {
	bst := tree.NewBinarySearchTree[string]()
	bst.Add("banana")
	bst.Add("apple")
	bst.Add("cherry")
	got := bst.String()
	want := "[apple banana cherry]"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
