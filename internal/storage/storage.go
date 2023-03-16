package storage

import (
	"fmt"
	"log"
	"sync"
)

type rep struct {
	URLMap map[int]string
	mutex  *sync.RWMutex
	id     int
}

func New() *rep {
	return &rep{
		URLMap: make(map[int]string),
		mutex:  &sync.RWMutex{},
		id:     0,
	}
}

func (r *rep) ReadURL(id int) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if str, ok := r.URLMap[id]; ok {
		log.Println("retrieved URL with ID:", id, "from URLMap", str)
		return str, nil
	}
	err := fmt.Errorf("URL with ID %d is not found in URLMap", id)
	return "", err
}

func (r *rep) WriteURL(url string) int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.id++
	r.URLMap[r.id] = url
	log.Println("added URL with ID", r.id, "to URLMap:", url)
	return r.id
}
