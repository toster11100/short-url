package storage

type Rep map[int]string

func New() Rep {
	rep := make(map[int]string)
	return rep
}

func (r Rep) ReadURL(id int) string {
	return r[id]
}

func (r Rep) WriteURL(url string, id int) {
	r[id] = url
}
