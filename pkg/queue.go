package painted

import "github.com/gammazero/deque"

type IoQueue struct {
	queue         deque.Deque
	index         int
	PrintCallback func(*Notification)
}

func (i *IoQueue) Next() {
	if i.index > 0 {
		i.index -= 1
	}
}

func (i *IoQueue) Prev() {
	if i.index+1 < i.queue.Len() {
		i.index += 1
	}
}

func (i *IoQueue) Push(n *Notification) {
	i.index = 0
	i.queue.PushFront(n)
}

func (i *IoQueue) Display() {
	n := i.queue.At(i.index).(*Notification)
	i.PrintCallback(n)
}
