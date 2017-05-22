package testutil

import (
	"errors"
	"sync"
)

type Body struct{}

var body []byte
var bodyLock sync.Mutex

func (b Body) Write(p []byte) (n int, err error) {
	body = make([]byte, len(p))
	bodyLock.Lock()
	n = copy(body, p)
	bodyLock.Unlock()
	return n, nil
}

func (b Body) Read(p []byte) (n int, err error) {
	if len(body) < 1 {
		return 0, errors.New("Nothing to write")
	}
	bodyLock.Lock()
	n = copy(p, body)
	bodyLock.Unlock()
	return n, nil
}
