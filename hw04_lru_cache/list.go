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
	return l.listLen
}

func (l *list) Front() *ListItem {
	return l.frontNode
}

func (l *list) Back() *ListItem {
	return l.backNode
}

// add a new node to the front of the list
func (l *list) PushFront(v interface{}) *ListItem {
	switch {
	case l.listLen == 0: // check if the list is empty
		l.frontNode = &ListItem{
			value: v,
			next:  nil,
			prev:  nil,
		}
		l.backNode = l.frontNode
	default:
		l.frontNode.prev = &ListItem{
			value: v,
			next:  l.frontNode,
			prev:  nil,
		}
		l.frontNode = l.frontNode.prev
	}
	l.listLen++
	return l.frontNode
}

// add a new node to the back of the list
func (l *list) PushBack(v interface{}) *ListItem {
	switch {
	case l.listLen == 0: // check if the list is empty
		l.backNode = &ListItem{
			value: v,
			next:  nil,
			prev:  nil,
		}
		l.frontNode = l.backNode
	default:
		l.backNode.next = &ListItem{
			value: v,
			next:  nil,
			prev:  l.backNode,
		}
		l.backNode = l.backNode.next
	}
	l.listLen++
	return l.backNode
}

// delete a specific node from the list
func (l *list) Remove(i *ListItem) {
	switch {
	case l.Len() == 1: // if we're deleting the only item in the list
		l.frontNode = nil
		l.backNode = nil
	case i.next == nil: // if we're deleting the Back node
		i.prev.next = nil
		l.backNode = i.prev
	case i.prev == nil: // if we're deleting the Front node
		i.next.prev = nil
		l.frontNode = i.next
	default: // If we're deleting a node from the middle
		i.next.prev = i.prev
		i.prev.next = i.next
	}
	l.listLen--
}

// moves a specific node to the front of the list
func (l *list) MoveToFront(i *ListItem) {
	if i == l.frontNode { // if trying to move the element that's already in front
		return
	}
	l.PushFront(i.value)
	l.Remove(i)
}

type ListItem struct {
	value   interface{}
	next    *ListItem
	prev    *ListItem
	itemKey Key // extra field for usage in cache
}

type list struct {
	frontNode *ListItem
	backNode  *ListItem
	listLen   int
}

func NewList() List {
	return new(list)
}
