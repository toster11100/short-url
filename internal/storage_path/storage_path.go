package storagepath

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type Rep struct {
	path    string
	mutex   *sync.RWMutex
	storage *storageJSON
}

type storageJSON struct {
	Max int       `json:"max"`
	URL []URLInfo `json:"url"`
}

type URLInfo struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func New(storagePath string) *Rep {
	return &Rep{
		path:    storagePath,
		mutex:   &sync.RWMutex{},
		storage: &storageJSON{Max: 1},
	}
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(path string) (*producer, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteFromStorage(event *storageJSON) error {
	if err := p.encoder.Encode(&event); err != nil {
		return err
	}
	return nil
}

func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(path string) (*consumer, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) ReadFromStorage() (*storageJSON, error) {
	event := &storageJSON{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func (r *Rep) ReadURL(id int) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	consumer, err := NewConsumer(r.path)
	if err != nil {
		return "", fmt.Errorf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	storage, err := consumer.ReadFromStorage()
	if err != nil {
		return "", nil
	}

	for _, urlInfo := range storage.URL {
		if urlInfo.ID == id {
			log.Printf("retrieved URL with ID %d from file: %s\n", id, urlInfo.URL)
			return urlInfo.URL, nil
		}
	}

	return "", fmt.Errorf("URL with ID %d is not found in file", id)
}

func (r *Rep) WriteURL(url string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	producer, err := NewProducer(r.path)
	if err != nil {
		return 0, fmt.Errorf("failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := NewConsumer(r.path)
	if err != nil {
		return 0, fmt.Errorf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	fileInfo, err := os.Stat(producer.file.Name())
	if err != nil {
		return 0, fmt.Errorf("failed to open file")
	}
	if fileInfo.Size() == 0 {
		r.storage.URL = append(r.storage.URL, URLInfo{ID: r.storage.Max, URL: url})
	} else {
		storage, err := consumer.ReadFromStorage()
		if err != nil {
			return 0, fmt.Errorf("failed to read from storage: %v", err)
		}
		r.storage.Max = storage.Max + 1
		r.storage.URL = append(storage.URL, URLInfo{ID: r.storage.Max, URL: url})
	}

	err = producer.WriteFromStorage(r.storage)
	if err != nil {
		return 0, fmt.Errorf("failed to write to storage: %v", err)
	}

	log.Println("added URL with ID", r.storage.Max, "to file:", url)
	return r.storage.Max, nil
}
