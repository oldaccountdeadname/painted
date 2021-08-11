package main

import "github.com/gammazero/deque"

type IoQueue struct {
	queue deque.Deque
	Model *Model
}

func (i *IoQueue) Push(n *Notification) {
	i.queue.PushFront(n)
}

func (i *IoQueue) Display() {
	n := i.queue.Front().(*Notification)
	i.Model.Io.Writef("%+v\n", n)
}
