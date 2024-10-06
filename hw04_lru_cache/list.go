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

func (l *list) PushFront(v interface{}) *ListItem {
	switch {
	case l.ListLen == 0:
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

func (l *list) PushBack(v interface{}) *ListItem {
	switch {
	case l.ListLen == 0:
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

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	l.ListLen--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
	l.ListLen--
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	FrontNode *ListItem
	BackNode  *ListItem
	ListLen   int
}

func NewList() List {
	return new(list)
}
