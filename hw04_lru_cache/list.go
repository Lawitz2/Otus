package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

func (l *list) Len() int {
	return l.ListLen
}

func (l *list) Front() *ListItem {
	return l.FrontNode
}

func (l *list) Back() *ListItem {
	return l.BackNode
}

// add a new node to the front of the list
func (l *list) PushFront(v interface{}) *ListItem {
	switch {
	case l.ListLen == 0: // check if the list is empty
		l.FrontNode = &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}
		l.BackNode = l.FrontNode
	default:
		l.FrontNode.Prev = &ListItem{
			Value: v,
			Next:  l.FrontNode,
			Prev:  nil,
		}
		l.FrontNode = l.FrontNode.Prev
	}
	l.ListLen++
	return l.FrontNode
}

// add a new node to the back of the list
func (l *list) PushBack(v interface{}) *ListItem {
	switch {
	case l.ListLen == 0: // check if the list is empty
		l.BackNode = &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}
		l.FrontNode = l.BackNode
	default:
		l.BackNode.Next = &ListItem{
			Value: v,
			Next:  nil,
			Prev:  l.BackNode,
		}
		l.BackNode = l.BackNode.Next
	}
	l.ListLen++
	return l.BackNode
}

// delete a specific node from the list
func (l *list) Remove(i *ListItem) {
	switch {
	case l.Len() == 1: // if we're deleting the only item in the list
		l.FrontNode = nil
		l.BackNode = nil
	case i.Next == nil: // if we're deleting the Back node
		i.Prev.Next = nil
		l.BackNode = i.Prev
	case i.Prev == nil: // if we're deleting the Front node
		i.Next.Prev = nil
		l.FrontNode = i.Next
	default: // If we're deleting a node from the middle
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.ListLen--
}

// moves a specific node to the front of the list
func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

type ListItem struct {
	Value   interface{}
	Next    *ListItem
	Prev    *ListItem
	ItemKey Key // extra field for usage in cache
}

type list struct {
	FrontNode *ListItem
	BackNode  *ListItem
	ListLen   int
}

func NewList() List {
	return new(list)
}
