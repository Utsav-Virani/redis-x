// -----------------------------------------------------
// REDIS-X
// © Utsav Virani
// Written by: (Utsav Virani)
// -----------------------------------------------------

package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type AOF struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

/*
- What happens here is that we first create the file if it doesn’t exist or open it if it does.
- Then, we create the bufio.Reader to read from the file.
- We start a goroutine to sync the AOF file to disk every 1 second while the server is running.
*/

func NewAOF(file string) (*AOF, error) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	aof := &AOF{
		file: f,
		rd:   bufio.NewReader(f),
	}
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()
	return aof, nil
}

func (aof *AOF) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *AOF) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Driver())
	if err != nil {
		return err
	}

	return nil
}
