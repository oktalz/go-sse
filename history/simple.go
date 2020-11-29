package history

import (
	"fmt"
)

var ErrMaxSizeZero = fmt.Errorf("maxSize is zero")
var ErrIDNotExists = fmt.Errorf("ID does not exists")

type Simple struct {
	MaxSize   int
	data      map[int]Message
	length    int
	nextID    int
	minUsedID int
}

func (m *Simple) Add(msg Message) (id int, err error) {
	if m.MaxSize == 0 {
		return -1, ErrMaxSizeZero
	}
	if m.data == nil {
		m.data = map[int]Message{}
	}
	if m.length == m.MaxSize {
		delete(m.data, m.minUsedID)
		m.minUsedID++
	} else {
		m.length++
	}
	m.data[m.nextID] = msg
	m.nextID++
	return m.nextID - 1, nil
}

// Get returns Message, use errors.Is(err, ErrIDNotExists)
func (m *Simple) Get(id int) (Message, error) {
	if id < m.minUsedID {
		return Message{}, fmt.Errorf("%w: %d", ErrIDNotExists, id)
	}
	if id >= m.nextID {
		return Message{}, fmt.Errorf("%w: %d", ErrIDNotExists, id)
	}
	msg, ok := m.data[id]
	if !ok {
		return Message{}, fmt.Errorf("%w: %d", ErrIDNotExists, id)
	}
	return msg, nil
}

// Last returns last Message, use errors.Is(err, ErrIDNotExists)
func (m *Simple) Last() (Message, error) {
	if m.length == 0 {
		return Message{}, ErrNoMessages
	}
	msg := m.data[m.nextID-1]
	return msg, nil
}
