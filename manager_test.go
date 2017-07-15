package queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type work int

func (w *work) Do(v interface{}) {
	fmt.Println(w, v)
	*w++
}

func TestQueuManager(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var w work
	q := NewQueue(len(items))
	ctx := context.Background()
	m := NewManager(ctx, q, &w, items...)

	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())

	<-m.End()

	if w != 10 {
		t.Error("not finish", w)
	}
}

type slow int

func (w *slow) Do(v interface{}) {
	time.Sleep(500 * time.Millisecond)
	*w++
	fmt.Println(w, v)
}

func TestQueuManagerTimeout(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var w slow
	q := NewQueue(len(items))
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	m := NewManager(ctx, q, &w, items...)

	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())

	<-m.End()

	if w == 10 {
		t.Error("not timeout", w)
	}
}