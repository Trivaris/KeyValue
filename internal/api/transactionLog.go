package api

import (
	"bufio"
	"os"
	"sync"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type TransactionLog struct {
	file              *os.File
	logSequenceNumber int
	mu                sync.Mutex
}

func NewTransactionLog(path string) (*TransactionLog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &TransactionLog{file: file}, nil
}

func (l *TransactionLog) Log(commandArray serialization.Array) {
	l.mu.Lock()
	defer l.mu.Unlock()

	r := bufio.NewWriter(l.file)
	r.Write(commandArray.Serialize())
}

func (l *TransactionLog) Close() {
	defer func() { _ = l.file.Close() }()
}
