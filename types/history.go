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
		st:   -1,
		ed:   -1,
	}
}

func (h *History) Push(val ResMessage) {
	newEd := (h.ed + 1) % h.len

	if h.st == -1 && h.ed == -1 {
		h.st = 0
		h.ed = 0
		h.data[h.ed] = val
	} else if newEd == h.st {
		h.st = (h.st + 1) % h.len
		h.ed = newEd
		h.data[h.ed] = val
	} else {
		h.ed = newEd
		h.data[h.ed] = val
	}
}

func (h *History) Iter() *HistoryIter {
	return &HistoryIter{
		data: h.data,
		len:  len(h.data),
		cur:  h.st,
		ed:   h.ed,
	}
}

type HistoryIter struct {
	data []ResMessage
	len  int
	cur  int
	ed   int
}

func (i *HistoryIter) Next() bool {
	if i.cur == -1 && i.ed == -1 {
		return false
	}
	return i.cur != (i.ed+1)%i.len
}

func (i *HistoryIter) Get() ResMessage {
	res := i.data[i.cur]
	i.cur = (i.cur + 1) % i.len
	return res
}
