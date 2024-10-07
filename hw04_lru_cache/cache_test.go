package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		c.Clear() // new test
		val, ok = c.Get("aaa")
		require.False(t, ok)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(4)
		c.Set("a", 2)
		c.Set("b", 4)
		c.Set("c", 6)
		c.Set("d", 8) // [d:8, c:6, b:4, a:2]

		val, ok := c.Get("a") // [a:2, d:8, c:6, b:4]
		require.Equal(t, 2, val)
		require.Equal(t, true, ok)

		val, ok = c.Get("d") // [d:8, a:2, c:6, b:4]
		require.Equal(t, 8, val)
		require.Equal(t, true, ok)

		c.Set("e", 10) // [e:10, d:8, a:2, c:6]
		c.Set("f", 12) // [f:12, e:10, d:8, a:2]

		val, ok = c.Get("c") // [f:12, e:10, d:8, a:2]
		require.Equal(t, false, ok)

		ok = c.Set("a", 14) // [a:14, f:12, e:10, d:8]
		require.Equal(t, true, ok)

		c.Get("d")     // [d:8, a:14, f:12, e:10]
		c.Get("e")     // [e:10, d:8, a:14, f:12]
		c.Set("g", 16) // [g:16, e:10, d:8, a:14]

		val, ok = c.Get("f") // [g:16, e:10, d:8, a:14]
		require.Equal(t, false, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
