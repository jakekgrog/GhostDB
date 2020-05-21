package linked_list

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	// Key of the key-value pair
	Key       string

	// Value of the key-value pair
	Value     string

	// TTL is the time-to-live for the key-value pair
	TTL       int64

	// CreatedAt is the time the key-value pair was entered
	// into the cache.
	CreatedAt int64

	// Prev points to the previous node in the doubly
	// linked list. Omit this from snapshot serialization.
	Prev      *Node `json:"-"`

	// Next points to the next node in the doubly linked
	// list. Omit this from snapshot serialization.
	Next      *Node `json:"-"`

	// Mux is a mutex lock.
	Mux       sync.Mutex
}

type List struct {
	// Head is the head node. It is a special case node.
	// It does not get populated and is a reference node 
	// for accessing the most recently used key-value pair.
	Head *Node `json:"-"`

	// Tail is the tail node. It is a special case node.
	// It does not get populated and is a reference node
	// for accessing the least recently used key-value pair.
	Tail *Node `json:"-"`

	// Size is the size of the list.
	Size int32
	Mux  sync.Mutex
}

// InitList initializes the doubly-linked list.
func InitList() *List {

	// Init the head node
	headNode := &Node{
		Key:       "",
		Value:     "",
		TTL:       -1,
		CreatedAt: time.Now().Unix(),
		Prev:      nil,
		Next:      nil,
	}

	// Init the tail node
	tailNode := &Node{
		Key:       "",
		Value:     "",
		TTL:       -1,
		CreatedAt: time.Now().Unix(),
		Prev:      nil,
		Next:      nil,
	}

	// Init the doubly-linked list
	list := &List{
		Head: headNode,
		Tail: tailNode,
		Size: int32(0),
	}

	// Set correct pointers for head and tail nodes.1
	list.Head.Next = list.Tail
	list.Tail.Prev = list.Head

	return list
}

// Insert will insert key-value pairs nodes into the doubly
// linked list.
func Insert(ll *List, key string, value string, ttl int64) (*Node, error) {
	// Lock access to the list
	ll.Mux.Lock()
	defer ll.Mux.Unlock()

	// Init the new node
	newNode := &Node{
		Key:       key,
		Value:     value,
		TTL:       ttl,
		CreatedAt: time.Now().Unix(),
		Prev:      nil,
		Next:      nil,
	}

	// Update the pointers of head and tail and set pointers
	// for the new node.
	newNode.Prev = ll.Head
	newNode.Next = ll.Head.Next
	ll.Head.Next = newNode      // Point Head to newNode
	newNode.Next.Prev = newNode // Point the old "Most Recent" to the new node

	// Atomically increment the size.
	atomic.AddInt32(&ll.Size, 1)

	return newNode, nil
}

// RemoveLast removes the least recently used item in the list.
func RemoveLast(ll *List) (*Node, error) {
	// Lock access
	ll.Mux.Lock()
	defer ll.Mux.Unlock()
	
	if ll.Size == 0 {
		return nil, errors.New("List is empty")
	} else {
		// Update reference pointers			
		nodeToRemove := ll.Tail.Prev

		nodeToRemove.Prev.Next = ll.Tail		
		ll.Tail.Prev = nodeToRemove.Prev
		
		atomic.AddInt32(&ll.Size, -1)

		return nodeToRemove, nil
	}
}

// RemoveNode removes a specific node from the list.
func RemoveNode(ll *List, node *Node) (*Node, error) {
	ll.Mux.Lock()
	if ll.Size == 0 {
		ll.Mux.Unlock()
		return nil, errors.New("List is empty")
	} else if ll.Size == 1 {
		ll.Mux.Unlock()
		returnNode, _ := RemoveLast(ll)
		return returnNode, nil
	}
	ll.Mux.Unlock()

	ll.Mux.Lock()
	defer ll.Mux.Unlock()

	prevNode := node.Prev
	nextNode := node.Next
	
	prevNode.Next = node.Next
	nextNode.Prev = node.Prev

	atomic.AddInt32(&ll.Size, -1)

	return node, nil
}

// Returns the last node in the list
func GetLastNode(ll *List) (*Node, error) {
	ll.Mux.Lock()
	if ll.Size == int32(0) {
		ll.Mux.Unlock()
		return nil, errors.New("List is empty")
	}
	ll.Mux.Unlock()

	nodeToGet := ll.Tail.Prev
	return nodeToGet, nil
}
