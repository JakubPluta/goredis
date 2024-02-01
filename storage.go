package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

type RedisFileStore struct {
	// Define the RedisFileStore struct
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewRedisFileStore(filePath string) (*RedisFileStore, error) {
	// What happens here is that we first create the file if it doesn’t exist or open it if it does.
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	rfs := &RedisFileStore{
		file: file,
		rd:   bufio.NewReader(file), // read from the file with a buffered reader
	}
	// start a goroutine to sync the file to disk every 1 second while the server is running.
	go func() {
		for {
			rfs.mu.Lock()
			rfs.file.Sync()
			rfs.mu.Unlock()
			time.Sleep(time.Second) // syncing every second ensures that the changes we made are always present on disk.
		}
	}()
	return rfs, nil

}

func (rfs *RedisFileStore) Write(m RedisMessage) error {
	rfs.mu.Lock()
	defer rfs.mu.Unlock()

	_, err := rfs.file.Write(m.MarshalMessage())
	if err != nil {
		return err
	}
	return nil

}

// Close closes the RedisFileStore.
//
// No parameters.
// Returns an error.
func (rfs *RedisFileStore) Close() error {
	rfs.mu.Lock()
	defer rfs.mu.Unlock()
	return rfs.file.Close()
}

func (rfs *RedisFileStore) Read(fn func(msg RedisMessage)) error {
	rfs.mu.Lock()
	defer rfs.mu.Unlock()

	rfs.file.Seek(0, io.SeekStart) // seek to the beginning of the file

	reader := NewResp(rfs.file)

	for {
		msg, err := reader.ReadMessage()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fn(msg)

	}
	return nil
}
