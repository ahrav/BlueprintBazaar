package chatroom

import (
	"container/heap"
	"errors"
	"sync"
)

type User struct {
	id           string
	messageCount int
	idx          int
}

type UserMaxHeap []*User

func (h UserMaxHeap) Len() int           { return len(h) }
func (h UserMaxHeap) Less(i, j int) bool { return h[i].messageCount > h[j].messageCount }
func (h UserMaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}

func (h *UserMaxHeap) Push(x any) {
	n := len(*h)
	user := x.(*User)
	user.idx = n
	*h = append(*h, user)
}

func (h *UserMaxHeap) Pop() any {
	old := *h
	n := len(old)
	user := old[n-1]
	old[n-1] = nil
	user.idx = -1
	*h = old[:n-1]
	return user
}

func (h *UserMaxHeap) Update(user *User, msgCnt int) {
	user.messageCount = msgCnt
	heap.Fix(h, user.idx)
}

type ChatRoom struct {
	id string

	mu    sync.RWMutex
	users map[string]*User

	userMaxHeap *UserMaxHeap
}

type ChatSystem struct {
	mu        sync.RWMutex
	chatRooms map[string]*ChatRoom
}

func (cr *ChatRoom) GetTopUser() (*User, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	usr, ok := cr.users[(*cr.userMaxHeap)[0].id]
	if !ok {
		return nil, errors.New("invalid userID, user not present in system")
	}
	return usr, nil
}

func (cs *ChatSystem) SendMessage(roomID, userID string) error {
	cs.mu.RLock()
	room, ok := cs.chatRooms[roomID]
	cs.mu.RUnlock()
	if !ok {
		return errors.New("room does not exist")
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if user, exists := room.users[userID]; exists {
		user.messageCount++
		room.userMaxHeap.Update(user, user.messageCount)
	} else {
		usr := &User{messageCount: 1, id: userID}
		heap.Push(room.userMaxHeap, usr)
	}

	return nil
}

func (cs *ChatSystem) GetGlobalTopUser() *User {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	var topUser *User
	for _, room := range cs.chatRooms {
		tUser, _ := room.GetTopUser()
		if tUser != nil && (topUser == nil || topUser.messageCount < tUser.messageCount) {
			topUser = tUser
		}
	}

	return topUser
}
