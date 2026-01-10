# Data Structures and Algorithms - Homework

A Go library for learning data structures through implementation.

## Philosophy

This library uses small, composable interfaces defined in `adt/adt.go`. 
Each method you implement satisfies one or more interfaces. 

Tests in `adt/prop/` validate your implementations against these contracts.

## Getting Started

```bash
# run all tests (most will fail initially)
go test ./...

# run tests for a specific package
go test ./dynamicarray/...
go test ./linkedlist/...
```

## Homework Order

Complete the data structures in this order. Each builds on concepts from the previous:

| #  | Package        | Structure                     | Key Concepts                   |
|----|----------------|-------------------------------|--------------------------------|
| 1  | `dynamicarray` | `DynamicArray[T]`             | Amortized growth, index access |
| 2  | `linkedlist`   | `SinglyLinkedList[T]`         | Node pointers, traversal       |
| 3  | `linkedlist`   | `DoublyLinkedList[T]`         | Bidirectional links            |
| 4  | `linkedlist`   | `CircularDoublyLinkedList[T]` | Circular references            |
| 5  | `stack`        | `Stack[T]`                    | LIFO, backend abstraction      |
| 6  | `queue`        | `Queue[T]`                    | FIFO, backend abstraction      |
| 7  | `hashmap`      | `HashMap[K,V]`                | Hashing, collision resolution  |
| 8  | `sets`         | `HashSet[T]`                  | Set operations, uniqueness     |
| 9  | `graph`        | `Graph[V]`                    | Adjacency list, traversal      |
| 10 | `tree`         | `BinarySearchTree[T]`         | Tree structure, ordering       |
| 11 | `bitsets`      | `BitSet`                      | Bit manipulation (optional)    |
