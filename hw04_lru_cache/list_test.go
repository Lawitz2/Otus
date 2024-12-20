package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().value)
		require.Equal(t, 70, l.Back().value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.next {
			elems = append(elems, i.value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("item removal", func(t *testing.T) {
		l := NewList()

		l.PushFront(2)                        // [2]
		require.Equal(t, l.Front(), l.Back()) // Front == Back on a single-item List

		l.Remove(l.Front()) // []

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		for i := 2; i < 11; i += 2 {
			l.PushFront(i)
		} // [10, 8, 6, 4, 2]

		require.Equal(t, 10, l.Front().value)
		require.Equal(t, 2, l.Back().value)

		for l.Len() > 1 {
			switch l.Len() % 2 {
			case 0:
				l.Remove(l.Back())
			case 1:
				l.Remove(l.Front())
			}
		} // [6]

		require.Equal(t, 1, l.Len())
		require.Equal(t, 6, l.Front().value)
	})
}
