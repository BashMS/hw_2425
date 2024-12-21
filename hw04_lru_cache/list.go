package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size int
	head *ListItem
	tail *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{
		Value: v,
		Next:  l.head,
		Prev:  nil,
	}
	if l.head != nil {
		l.head.Prev = &item
	}
	if l.tail == nil {
		l.tail = &item
	}
	l.head = &item
	l.size++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
	}
	if l.tail != nil {
		l.tail.Next = &item
	}
	if l.head == nil {
		l.head = &item
	}
	l.tail = &item
	l.size++
	return &item
}

func (l *list) Remove(i *ListItem) {
	if l.size == 0 {
		return
	}
	if i == l.head {
		l.head = l.head.Next
		l.head.Prev = nil
		l.size--
		return
	}
	if i == l.tail {
		l.tail = l.tail.Prev
		l.tail.Next = nil
		l.size--
		return
	}

	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		l.size--
		return
	}
}

func (l *list) PrintAll() {
	if l.size == 0 {
		return
	}
	tmp := l.head
	for tmp != l.tail {
		fmt.Println(tmp.Value)
		tmp = tmp.Next
	}
	fmt.Println(tmp.Value)
}

func (l *list) MoveToFront(i *ListItem) {
	if l.size == 0 {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}
