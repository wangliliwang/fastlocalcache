package fastcache_lru

import (
	"fmt"
)

type linkedListNode[T any] struct {
	prev, next *linkedListNode[T]
	value      T
	key        string
}

type linkedList[T any] struct {
	head, tail *linkedListNode[T]
}

func newLinkedList[T any]() *linkedList[T] {
	return &linkedList[T]{
		head: nil,
		tail: nil,
	}
}

func (l *linkedList[T]) moveToHead(nd *linkedListNode[T]) {
	// new inserted
	if nd.prev == nil && nd.next == nil {
		if l.head == nil {
			l.head, l.tail = nd, nd
		} else {
			nd.next = l.head
			l.head.prev = nd
			l.head = nd
		}
		return
	}
	// nd is head
	if nd == l.head {
		return
	}
	// nd is tail
	if nd == l.tail {
		l.tail = l.tail.prev
		l.tail.next = nil
		nd.prev = nil
		nd.next = l.head
		l.head.prev = nd
		l.head = nd
		return
	}
	// nd is mid
	nd.prev.next = nd.next
	nd.next.prev = nd.prev
	nd.prev = nil
	nd.next = l.head
	l.head.prev = nd
	l.head = nd
}

func (l *linkedList[T]) delTail() *linkedListNode[T] {
	r := l.tail
	l.del(l.tail)
	return r
}

func (l *linkedList[T]) del(nd *linkedListNode[T]) {
	if nd == nil {
		return
	}
	if nd == l.head {
		if l.head == l.tail {
			l.head = nil
			l.tail = nil
		} else {
			l.head = l.head.next
			l.head.prev = nil
			nd.next = nil
		}
		return
	}
	if nd == l.tail {
		l.tail = l.tail.prev
		l.tail.next = nil
		nd.prev = nil
		return
	}
	nd.prev.next = nd.next
	nd.next.prev = nd.prev
	nd.next = nil
	nd.prev = nil
}

func (l *linkedList[T]) String() string {
	headS, tailS := "nil", "nil"
	if l.head != nil {
		headS = fmt.Sprintf("%+v", l.head.value)
		tailS = fmt.Sprintf("%+v", l.tail.value)
	}
	r := fmt.Sprintf("linkedList: \nhead=[%+v], tail=[%+v]\n", headS, tailS)
	for cur := l.head; cur != nil; cur = cur.next {
		r += fmt.Sprintf(" [%+v] ->", cur.value)
	}
	r += "\n"

	return r
}
