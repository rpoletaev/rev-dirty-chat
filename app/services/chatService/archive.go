package chatService

import (
	"container/list"
)

type archive struct {
	lst     *list.List
	maxSize int
}

func CreateArchive(sz int) *archive {
	return &archive{
		lst:     list.New(),
		maxSize: sz,
	}
}

func (arc *archive) Add(message interface{}) {
	if arc.lst.Len() >= arc.maxSize {
		arc.lst.Remove(arc.lst.Front())
	}

	arc.lst.PushBack(message)
}

func (arc *archive) Each(fn func(message interface{})) {
	for msg := arc.lst.Front(); msg != nil; msg = msg.Next() {
		fn(msg.Value)
	}
}
