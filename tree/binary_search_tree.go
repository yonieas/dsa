// Package tree provides tree data structure implementations.
//
// # What is a Tree?
//
// A tree is a hierarchical data structure with a root node and subtrees of
// children. File systems, organization charts, HTML DOM, and taxonomy
// classifications are all trees.
//
// # Terminology
//
//	Root:     The topmost node (has no parent)
//	Parent:   A node that has children
//	Child:    A node connected below its parent
//	Leaf:     A node with no children
//	Height:   Longest path from root to any leaf
//	Depth:    Distance from root to a specific node
//
// # Binary Search Tree (BST)
//
// A BST is a binary tree (each node has at most two children) with a special
// ordering property: all values in the left subtree are less than the node,
// and all values in the right subtree are greater. This ordering enables
// efficient search, insert, and delete.
//
// Finding a value in a BST works like binary search: compare with the current
// node, go left if smaller, right if larger. Each comparison eliminates half
// the remaining nodes, giving O(log n) average time.
//
// # Traversals
//
//	In-Order:   Left, Root, Right (yields sorted order for BST)
//	Pre-Order:  Root, Left, Right (useful for copying trees)
//	Post-Order: Left, Right, Root (useful for deleting trees)
//
// # Complexity
//
//	Search/Insert/Delete: O(log n) average, O(n) worst case
//	Traversal:            O(n)
//	Space:                O(n)
//
// The O(n) worst case happens when the tree becomes a linked list (inserting
// sorted data). Self-balancing trees like AVL or Red-Black trees guarantee
// O(log n) by keeping the tree balanced.
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 12.
// Sedgewick "Algorithms", Section 3.2.
// https://en.wikipedia.org/wiki/Binary_search_tree
package tree

import (
	"cmp"

	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/sequence"
)

// Node is a node in a binary search tree.
// Each node holds data and pointers to left and right children.
//
//	       ┌─────────┐
//	       │  data   │
//	       └────┬────┘
//	      ╱           ╲
//	   left           right
//	(< data)         (> data)
type Node[E cmp.Ordered] struct {
	data  E
	left  *Node[E]
	right *Node[E]
}

// BinarySearchTree is a binary tree with the BST ordering property.
//
// For any node N:
//   - All nodes in left subtree have values < N.data
//   - All nodes in right subtree have values > N.data
//
// Example BST with values 5, 3, 7, 1, 4, 6, 8:
//
//	     5          ← root
//	   /   \
//	  3     7
//	 / \   / \
//	1   4 6   8
//
// This ordering means:
//   - In-order traversal yields sorted values: 1, 3, 4, 5, 6, 7, 8
//   - Search can eliminate half the tree at each step
//
// Note: This is an unbalanced BST. Inserting sorted values creates
// a degenerate tree (essentially a linked list) with O(n) operations.
// Use self-balancing trees (AVL, Red-Black) for guaranteed O(log n).
type BinarySearchTree[E cmp.Ordered] struct {
	root *Node[E]
	size int
}

// NewBinarySearchTree creates an empty binary search tree.
//
//	root = nil
//	size = 0
func NewBinarySearchTree[E cmp.Ordered]() *BinarySearchTree[E] {
	return &BinarySearchTree[E]{}
}

// Size returns the number of nodes in the tree.
//
//	            5
//	          /   \
//	         3     7
//	        / \
//	       1   4
//
//	Size() → 5
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (t *BinarySearchTree[E]) Size() int {
	return t.size
}

// Empty returns true if the tree has no nodes.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (t *BinarySearchTree[E]) Empty() bool {
	return t.size == 0
}

// Add inserts a value into the BST, maintaining the ordering property.
//
//	Before Add(4):
//
//	        5
//	      /   \
//	     3     7
//	    /
//	   1
//
//	After Add(4):
//
//	        5
//	      /   \
//	     3     7
//	    / \
//	   1   4  ← inserted here (4 > 3, so goes right of 3)
//
// Insertion process:
//  1. Start at root
//  2. If value < current node, go left
//  3. If value > current node, go right
//  4. Repeat until finding nil, insert there
//
// Duplicates are ignored (value already exists → no change).
//
// complexity:
//   - time : O(h) where h is the height of the tree
//   - space: O(h) for the recursive call stack
//     For balanced tree: h = O(log n)
//     For degenerate tree: h = O(n)
//
// SCORE: 20
func (t *BinarySearchTree[E]) Add(value E) {
	// hint: use recursive helper function:
	//       1) if node == nil, create new node, return (node, true)
	//       2) if value < node.data, recurse left
	//       3) if value > node.data, recurse right
	//       4) if equal, do nothing (duplicate)
	//       5) increment size if inserted
	var inserted bool
	root, inserted := t.addHelper(t.root, value)
	if inserted {
		t.size++
	}
	t.root = root
}

// Recursive helper for `Add` method
func (t *BinarySearchTree[E]) addHelper(node *Node[E], value E) (*Node[E], bool) {
	if node == nil {
		return &Node[E]{data: value}, true
	}
	if value < node.data {
		newLeft, inserted := t.addHelper(node.left, value)
		node.left = newLeft
		return node, inserted
	} else if value > node.data {
		newRight, inserted := t.addHelper(node.right, value)
		node.right = newRight
		return node, inserted
	} else {
		return node, false
	}
}

// Exists checks if a value exists in the BST.
//
//	        5
//	      /   \
//	     3     7
//	    / \   / \
//	   1   4 6   8
//
//	Exists(4):
//	  - Start at 5: 4 < 5, go left
//	  - At 3: 4 > 3, go right
//	  - At 4: found! → return true
//
//	Exists(9):
//	  - Start at 5: 9 > 5, go right
//	  - At 7: 9 > 7, go right
//	  - At 8: 9 > 8, go right
//	  - At nil: not found → return false
//
// complexity:
//   - time : O(h) where h is the height of the tree
//   - space: O(h) for the recursive call stack
//
// SCORE: 15
func (t *BinarySearchTree[E]) Exists(value E) bool {
	// hint: use recursive helper function:
	//       1) if node == nil, return false
	//       2) if value < node.data, search left
	//       3) if value > node.data, search right
	//       4) if equal, return true
	return t.existsHelper(t.root, value)
}

// Recursive helper for `Exist` method
func (t *BinarySearchTree[E]) existsHelper(node *Node[E], value E) bool {
	if node == nil {
		return false
	}
	if value < node.data {
		return t.existsHelper(node.left, value)
	} else if value > node.data {
		return t.existsHelper(node.right, value)
	} else {
		return true
	}
}

// Del removes a value from the BST, maintaining the ordering property.
//
// Deletion has three cases:
//
// Case 1: Leaf node (no children)
// Simply remove the node.
//
//	Before Del(4):           After Del(4):
//
//	    5                        5
//	   / \                      / \
//	  3   7                    3   7
//	   \
//	    4 ← remove
//
// Case 2: Node with one child
// Replace node with its child.
//
//	Before Del(3):           After Del(3):
//
//	    5                        5
//	   / \                      / \
//	  3   7                    1   7
//	 /
//	1   ← child moves up
//
// Case 3: Node with two children
// Replace with in-order successor (smallest in right subtree).
//
//	Before Del(5):           After Del(5):
//
//	    5 ← delete               6 ← successor
//	   / \                      / \
//	  3   7                    3   7
//	     /                          \
//	    6 ← successor                8
//	     \
//	      8
//
// complexity:
//   - time : O(h) where h is the height of the tree
//   - space: O(h) for the recursive call stack
//
// SCORE: 20
func (t *BinarySearchTree[E]) Del(value E) {
	// hint: use recursive helper function:
	//       1) if node == nil, return (nil, false)
	//       2) if value < node.data, recurse left
	//       3) if value > node.data, recurse right
	//       4) if equal: 3 cases
	//          - no left child: return right child
	//          - no right child: return left child
	//          - both children: find min in right subtree (successor),
	//            copy value, delete successor from right subtree
	//       5) decrement size if deleted
	var deleted bool
	t.root, deleted = t.deleteHelper(t.root, value)
	if deleted {
		t.size--
	}
}

// Recursive helper for `Del` method
func (t *BinarySearchTree[E]) deleteHelper(node *Node[E], value E) (*Node[E], bool) {
	if node == nil {
		return nil, false
	}
	if value < node.data {
		newLeft, deleted := t.deleteHelper(node.left, value)
		node.left = newLeft
		return node, deleted
	} else if value > node.data {
		newRight, deleted := t.deleteHelper(node.right, value)
		node.right = newRight
		return node, deleted
	} else {
		// No left child
		if node.left == nil {
			return node.right, true
		}
		// No right child
		if node.right == nil {
			return node.left, true
		}
		// Have two children
		successor, _ := t.Min()
		node.data = successor
		node.right, _ = t.deleteHelper(node.right, successor)
		return node, true
	}
}

// Min returns the smallest value in the tree.
//
//	        5
//	      /   \
//	     3     7
//	    /
//	   1  ← minimum (leftmost node)
//
//	Min() → (1, true)
//
// The minimum is always the leftmost node in the tree.
//
// complexity:
//   - time : O(h) where h is the height of the tree
//   - space: O(1)
//
// SCORE: 10
func (t *BinarySearchTree[E]) Min() (E, bool) {
	// hint: 1) if root == nil, return (zero, false)
	//       2) traverse left until node.left == nil
	//       3) return (node.data, true)
	if t.root == nil {
		return generics.ZeroValue[E](), false
	}
	node := t.root
	for node.left != nil {
		node = node.left
	}
	return node.data, true
}

// Max returns the largest value in the tree.
//
//	        5
//	      /   \
//	     3     7
//	            \
//	             8  ← maximum (rightmost node)
//
//	Max() → (8, true)
//
// The maximum is always the rightmost node in the tree.
//
// complexity:
//   - time : O(h) where h is the height of the tree
//   - space: O(1)
//
// SCORE: 10
func (t *BinarySearchTree[E]) Max() (E, bool) {
	// hint: 1) if root == nil, return (zero, false)
	//       2) traverse right until node.right == nil
	//       3) return (node.data, true)
	if t.root == nil {
		return generics.ZeroValue[E](), false
	}
	node := t.root
	for node.right != nil {
		node = node.right
	}
	return node.data, true
}

// InOrder traverses the tree in sorted order (left → root → right).
//
//	        5
//	      /   \
//	     3     7
//	    / \   / \
//	   1   4 6   8
//
//	InOrder visits: 1, 3, 4, 5, 6, 7, 8  (sorted!)
//
//	Traversal order:
//	  1. Go to leftmost node (1)
//	  2. Visit 1
//	  3. Go up to 3
//	  4. Visit 3
//	  5. Go to right child (4)
//	  6. Visit 4
//	  ... and so on
//
// This is a key property of BST: in-order traversal yields sorted values.
//
// complexity:
//   - time : O(n) visits every node exactly once
//   - space: O(h) for the recursive call stack
//
// SCORE: 10
func (t *BinarySearchTree[E]) InOrder(visit func(E) bool) {
	// hint: use recursive helper: inOrder(node, visit)
	//       1) if node == nil, return true
	//       2) if !inOrder(node.left, visit), return false
	//       3) if !visit(node.data), return false
	//       4) return inOrder(node.right, visit)
	t.inOrderHelper(t.root, visit)
}

// Recursive helper for `InOrder` method
func (t *BinarySearchTree[E]) inOrderHelper(node *Node[E], visit func(E) bool) bool {
	if node == nil {
		return true
	}
	// Left traversal
	if !t.inOrderHelper(node.left, visit) {
		return false
	}
	// Visit root
	if !visit(node.data) {
		return false
	}
	// Right traversal
	return t.inOrderHelper(node.right, visit)
}

// Iter is an alias for InOrder, satisfying adt.Iterator.
// Iterating a BST yields values in sorted order.
func (t *BinarySearchTree[E]) Iter(visit func(E) bool) {
	t.InOrder(visit)
}

// String returns the string representation of the tree (in-order).
//
//	        5
//	      /   \
//	     3     7
//	    / \
//	   1   4
//
//	String() → "[1 3 4 5 7]"
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (t *BinarySearchTree[E]) String() string {
	return sequence.String(t.Iter)
}

// IterBackward traverses the tree in reverse in-order (right → root → left).
// This visits elements in descending order.
//
//	        5
//	      /   \
//	     3     7
//	    / \
//	   1   4
//
//	IterBackward visits: 7, 5, 4, 3, 1  (descending order)
//
// complexity:
//   - time : O(n)
//   - space: O(h)
//
// SCORE: 5
func (t *BinarySearchTree[E]) IterBackward(visit func(E) bool) {
	// hint: reverse of InOrder: right → root → left
	//       1) recurse right first
	//       2) visit node
	//       3) recurse left
	t.reverseInOrderHelper(t.root, visit)
}

// Recursive helper for `IterBackward` method
func (t *BinarySearchTree[E]) reverseInOrderHelper(node *Node[E], visit func(E) bool) bool {
	if node == nil {
		return true
	}
	// Right traversal
	if !t.reverseInOrderHelper(node.right, visit) {
		return false
	}
	// Visit root
	if !visit(node.data) {
		return false
	}
	// Left traversal
	return t.reverseInOrderHelper(node.left, visit)
}

// PreOrder traverses the tree in pre-order (root → left → right).
//
//	        5
//	      /   \
//	     3     7
//	    / \   / \
//	   1   4 6   8
//
//	PreOrder visits: 5, 3, 1, 4, 7, 6, 8  (root first)
//
//	Traversal order:
//	  1. Visit root (5)
//	  2. Traverse left subtree (3, 1, 4)
//	  3. Traverse right subtree (7, 6, 8)
//
// Use cases:
//   - Creating a copy of the tree
//   - Serializing/deserializing the tree
//   - Prefix expression evaluation
//
// complexity:
//   - time : O(n)
//   - space: O(h)
//
// SCORE: 5
func (t *BinarySearchTree[E]) PreOrder(visit func(E) bool) {
	// hint: root → left → right
	//       1) visit node first
	//       2) recurse left
	//       3) recurse right
	t.preOrderHelper(t.root, visit)
}

// Recursive helper for `PreOrder` method
func (t *BinarySearchTree[E]) preOrderHelper(node *Node[E], visit func(E) bool) bool {
	if node == nil {
		return true
	}
	// Visit root
	if !visit(node.data) {
		return false
	}
	// Left traversal
	if !t.preOrderHelper(node.left, visit) {
		return false
	}
	// Right traversal
	return t.preOrderHelper(node.right, visit)
}

// PostOrder traverses the tree in post-order (left → right → root).
//
//	        5
//	      /   \
//	     3     7
//	    / \   / \
//	   1   4 6   8
//
//	PostOrder visits: 1, 4, 3, 6, 8, 7, 5  (root last)
//
//	Traversal order:
//	  1. Traverse left subtree (1, 4, 3)
//	  2. Traverse right subtree (6, 8, 7)
//	  3. Visit root (5)
//
// Use cases:
//   - Deleting the tree (delete children before parent)
//   - Evaluating expression trees
//   - Calculating directory sizes
//
// complexity:
//   - time : O(n)
//   - space: O(h)
//
// SCORE: 5
func (t *BinarySearchTree[E]) PostOrder(visit func(E) bool) {
	// hint: left → right → root
	//       1) recurse left
	//       2) recurse right
	//       3) visit node last
	t.postOrderHelper(t.root, visit)
}

// Recursive function for `PostOrder` method
func (t *BinarySearchTree[E]) postOrderHelper(node *Node[E], visit func(E) bool) bool {
	if node == nil {
		return true
	}
	// Left traversal
	if !t.postOrderHelper(node.left, visit) {
		return false
	}
	// Right traversal
	if !t.postOrderHelper(node.right, visit) {
		return false
	}
	// Visit root
	if !visit(node.data) {
		return false
	}
	return true
}
