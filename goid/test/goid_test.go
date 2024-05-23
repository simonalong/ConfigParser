package test

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/simonalong/gole/goid"
	"github.com/simonalong/gole/test"
	"testing"
	"time"
)

func TestGoid(t *testing.T) {
	t.Log(goid.Goid())
}

func TestGet(t *testing.T) {
	ch := make(chan *string, 100)
	for i := 0; i < cap(ch); i++ {
		go func(i int) {
			gid := goid.NativeGoid()
			expected := goid.Goid()
			t.Logf("%d: %d", gid, expected)
			if gid == expected {
				ch <- nil
				return
			}
			s := fmt.Sprintf("Expected %d, but got %d", expected, gid)
			ch <- &s
		}(i)
	}

	for i := 0; i < cap(ch); i++ {
		val := <-ch
		if val != nil {
			t.Fatal(*val)
		}
	}
}

func TestAllGoid(t *testing.T) {
	const num = 10
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}
	time.Sleep(time.Millisecond)

	ids := goid.AllGoids()
	t.Log("all gids: ", len(ids), ids)
}

func TestGoStorage(t *testing.T) {
	var variable = "hello world"
	stg := goid.NewLocalStorage()
	stg.Set(variable)
	goid.Go(func() {
		v := stg.Get()
		test.True(t, v != nil && v.(string) == variable)
	})
	time.Sleep(time.Millisecond)
	stg.Clear()
}

// BenchmarkGoid-12    	278801190	         4.586 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGoid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = goid.Goid()
	}
}

// BenchmarkAllGoid-12    	 5949680	       228.3 ns/op	     896 B/op	       1 allocs/op
func BenchmarkAllGoid(b *testing.B) {
	const num = 16
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = goid.AllGoids()
	}
}

func BenchmarkNativeGoid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = goid.NativeGoid()
	}
}

var h1 goid.LocalStorage
var h2 goid.LocalStorage
var h3 goid.LocalStorage
var h4 goid.LocalStorage

func init() {
	h1 = goid.NewLocalStorage()
	h2 = goid.NewLocalStorage()
	h3 = goid.NewLocalStorage()
	h4 = goid.NewLocalStorage()
}

func TestGoidGet(t *testing.T) {
	h1.Set("11")
	h2.Set("12")
	h3.Set("13")
	h4.Set("14")
	h1V := h1.Get()
	h2V := h2.Get()
	h3V := h3.Get()
	h4V := h4.Get()

	assert.Equal(t, "11", h1V)
	assert.Equal(t, "12", h2V)
	assert.Equal(t, "13", h3V)
	assert.Equal(t, "14", h4V)
}

func TestGoidClean(t *testing.T) {
	st := goid.NewLocalStorage()
	st.Set("vv")
	goid.Go(func() {
		fmt.Println("sleep")
		time.Sleep(2 * time.Second)
		fmt.Println("wake")

		fmt.Println(st.Get())
	})
	st.Del()
	fmt.Println("clean")
	time.Sleep(4 * time.Second)
	fmt.Println("end")
}

//func TestGoidChange(t *testing.T) {
//	st := goid.NewLocalStorage()
//	st.Set("vv")
//	goid.Go(func() {
//		fmt.Println("sleep")
//		time.Sleep(2*time.Second)
//		fmt.Println("wake")
//
//		fmt.Println(st.Get())
//	})
//	st.Set("vv0")
//	fmt.Println("clean")
//	time.Sleep(4*time.Second)
//	fmt.Println("end")
//}
//
//
//func TestGoidChange2(t *testing.T) {
//	st := goid.NewLocalStorage()
//	st.Set("vv")
//	goid.Go(func() {
//		fmt.Println("sleep")
//		time.Sleep(2*time.Second)
//		fmt.Println("setValue")
//		st.Set("vv-children")
//	})
//	time.Sleep(3*time.Second)
//	fmt.Println("parent-wake")
//	fmt.Println(st.Get())
//}
//
//func TestGoidChange3(t *testing.T) {
//	st := goid.NewLocalStorage()
//	st.Set("vv")
//	goid.Go(func() {
//		fmt.Println("sleep")
//		time.Sleep(2*time.Second)
//		fmt.Println("wake")
//		fmt.Println(st.Get())
//	})
//	st.Del()
//	fmt.Println("clean")
//	time.Sleep(3*time.Second)
//	fmt.Println("parent-wake")
//	fmt.Println(st.Get())
//}
