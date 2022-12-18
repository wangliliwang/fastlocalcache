package fastcache_lru

import (
	"fmt"
	"testing"
)

func TestLinkedList(t *testing.T) {
	// new
	c := newLinkedList[int]()
	fmt.Println("init: ", c)

	nd1 := &linkedListNode[int]{value: 1}
	c.moveToHead(nd1)
	fmt.Println("insert 1: ", c)

	nd2 := &linkedListNode[int]{value: 2}
	c.moveToHead(nd2)
	fmt.Println("insert 2: ", c)

	c.moveToHead(nd1)
	fmt.Println("move 1 to head: ", c)

	nd3 := &linkedListNode[int]{value: 3}
	c.moveToHead(nd3)
	fmt.Println("insert 3: ", c)

	nd4 := &linkedListNode[int]{value: 4}
	c.moveToHead(nd4)
	fmt.Println("insert 4: ", c)

	c.del(nd3)
	fmt.Println("del 3: ", c)

	c.del(nd2)
	fmt.Println("del 2: ", c)

	c.del(nd4)
	fmt.Println("del 4: ", c)

	c.delTail()
	fmt.Println("del tail: ", c)

	c.delTail()
	fmt.Println("del tail(nothing): ", c)
}
