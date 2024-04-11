package types

type History struct {
	data []ResMessage
	len  int
	st   int
	ed   int
}

func NewHistory(len int) *History {
	return &History{
		data: make([]ResMessage, len),
		len:  len,
		st:   0,
		ed:   1,
	}
}

func (h *History) Push(val ResMessage) {
	if h.ed == h.st {
		h.st++
	}
	h.data[h.ed] = val
	h.ed = (h.ed + 1) % h.len
}

func (h *History) Iter() *historyIter {
	return &historyIter{
		data: h.data,
		len:  len(h.data),
		cur:  h.st,
		ed:   h.ed,
	}
}

type historyIter struct {
	data []ResMessage
	len  int
	cur  int
	ed   int
}

func (i *historyIter) Next() bool {
	return (i.cur+1)%i.len != i.ed
}

func (i *historyIter) Get() ResMessage {
	i.cur = (i.cur + 1) % i.len
	return i.data[i.cur]
}
