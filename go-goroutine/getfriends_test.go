package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

//https://www.fullstory.com/blog/why-errgroup-withcontext-in-golang-server-handlers/

const (
	nWorkers = 10
)

var (
	users = []string{"alice", "bob", "charlie", "dave", "esther", "frank"}
)

type User struct {
	Id   int64
	Name string
}

// Mocks out an iterator to a remote datastore query.
type friendIterator struct {
	cursor int
}

func (it *friendIterator) Next(ctx context.Context) (int64, error) {
	if it.cursor >= len(users) {
		return 0, io.EOF
	}

	// Pretend each item takes 10ms to come back over the network.
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-time.After(10 * time.Millisecond):
		r := int64(it.cursor)
		it.cursor++
		log.Printf("found friend id: %d", r)
		return r, nil
	}
}

// Mocks out a remote friend list lookup.
func GetFriendIds(user int64) *friendIterator {
	return &friendIterator{}
}

// Mocks out a user profile lookup.
func GetUserProfile(ctx context.Context, id int64) (*User, error) {
	// Pretend the lookup takes 100ms to come back over the network.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(100 * time.Millisecond):
		if id < 0 || id >= int64(len(users)) {
			return nil, fmt.Errorf("unknown user: %d", id)
		}
		log.Printf("found user profile: %d", id)
		return &User{Id: id, Name: users[id]}, nil
	}
}

func GetFriends_Serial(ctx context.Context, user int64) (map[string]*User, error) {
	// Produce
	var friendIds []int64
	for it := GetFriendIds(user); ; {
		if id, err := it.Next(ctx); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("GetFriendIds %d: %s", user, err)
		} else {
			friendIds = append(friendIds, id)
		}
	}

	// Map
	ret := map[string]*User{}
	for _, friendId := range friendIds {
		if friend, err := GetUserProfile(ctx, friendId); err != nil {
			return nil, fmt.Errorf("GetUserProfile %d: %s", friendId, err)
		} else {
			ret[friend.Name] = friend
		}
	}
	return ret, nil
}

func GetFriends_Parallel(ctx context.Context, user int64) (map[string]*User, error) {
	friendIds := make(chan int64)

	// Produce
	go func() {
		defer close(friendIds)
		for it := GetFriendIds(user); ; {
			if id, err := it.Next(ctx); err != nil {
				if err == io.EOF {
					break
				}
				// What to do here?
				log.Fatalf("GetFriendIds %d: %s", user, err)
			} else {
				friendIds <- id
			}
		}
	}()

	friends := make(chan *User)

	// Map
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func() {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(friends)
				}
			}()

			for id := range friendIds { //friendIds 채널에 쓰여지는 경우에 처리됨
				if friend, err := GetUserProfile(ctx, id); err != nil {
					// What to do here?
					log.Fatalf("GetUserProfile %d: %s", user, err)
				} else {
					friends <- friend
				}
			}
		}()
	}

	// Reduce
	ret := map[string]*User{}
	for friend := range friends { //friends 채널에 쓰여지는 경우에
		ret[friend.Name] = friend
	}

	return ret, nil
}

func GetFriends_ErrGroup(ctx context.Context, user int64) (map[string]*User, error) {
	g, ctx := errgroup.WithContext(ctx)
	friendIds := make(chan int64)

	// Produce
	g.Go(func() error {
		defer close(friendIds)
		for it := GetFriendIds(user); ; {
			if id, err := it.Next(ctx); err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("GetFriendIds %d: %s", user, err)
			} else {
				friendIds <- id
			}
		}
	})

	friends := make(chan *User)

	// Map
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(friends)
				}
			}()

			for id := range friendIds {
				if friend, err := GetUserProfile(ctx, id); err != nil {
					return fmt.Errorf("GetUserProfile %d: %s", user, err)
				} else {
					friends <- friend
				}
			}
			return nil
		})
	}

	// Reduce
	ret := map[string]*User{}
	g.Go(func() error {
		for friend := range friends {
			ret[friend.Name] = friend
		}
		return nil
	})

	// Return the final result, and the error result from the subtask group.
	return ret, g.Wait()
}

func GetFriends_Selects(ctx context.Context, user int64) (map[string]*User, error) {
	g, ctx := errgroup.WithContext(ctx)
	friendIds := make(chan int64)

	// Produce
	g.Go(func() error {
		defer close(friendIds)
		for it := GetFriendIds(user); ; {
			if id, err := it.Next(ctx); err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("GetFriendIds %d: %s", user, err)
			} else {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case friendIds <- id:
				}
			}
		}
	})

	friends := make(chan *User)

	// Map
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(friends)
				}
			}()

			for id := range friendIds {
				if friend, err := GetUserProfile(ctx, id); err != nil {
					return fmt.Errorf("GetUserProfile %d: %s", user, err)
				} else {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case friends <- friend:
					}
				}
			}
			return nil
		})
	}

	// Reduce
	ret := map[string]*User{}
	g.Go(func() error {
		for friend := range friends {
			ret[friend.Name] = friend
		}
		return nil
	})

	// Return the final result, and the error result from the subtask group.
	return ret, g.Wait()
}

func TestGetFriends(t *testing.T) {
	type getFriendFunc func(ctx context.Context, user int64) (map[string]*User, error)
	type tc struct {
		name string
		fn   getFriendFunc
	}

	for _, tc := range []tc{
		{"GetFriends_Serial", GetFriends_Serial},
		{"GetFriends_Parallel", GetFriends_Parallel},
		{"GetFriends_ErrGroup", GetFriends_ErrGroup},
		{"GetFriends_Selects", GetFriends_Selects},
	} {
		log.Printf("%s", tc.name)

		ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		start := time.Now()
		rsp, err := tc.fn(ctx, 0)
		if err != nil {
			log.Fatalf("error: %s", err)
		} else {
			log.Printf("finished in %s: %s", time.Since(start), jsonString(rsp))
		}
		log.Println()
		cancel()
	}
}

func jsonString(v interface{}) string {
	if ret, err := json.Marshal(v); err != nil {
		return err.Error()
	} else {
		return string(ret)
	}
}
