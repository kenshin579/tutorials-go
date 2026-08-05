// Package singleflight provides an example of suppressing duplicate function
// calls with golang.org/x/sync/singleflight.
package singleflight

import (
	"fmt"

	"golang.org/x/sync/singleflight"
)

// User is the value returned by the example data source.
type User struct {
	ID   string
	Name string
}

// FetchUserFunc represents a slow database or external API call.
type FetchUserFunc func(userID string) (User, error)

// UserLoader combines concurrent requests for the same user ID.
//
// UserLoader does not cache completed results. It only suppresses duplicate
// calls while a call with the same key is still in flight.
type UserLoader struct {
	group singleflight.Group
	fetch FetchUserFunc
}

// NewUserLoader creates a UserLoader backed by fetch.
func NewUserLoader(fetch FetchUserFunc) *UserLoader {
	return &UserLoader{fetch: fetch}
}

// Load returns a user and reports whether the result was shared by multiple
// callers. Concurrent calls using different user IDs run independently.
func (l *UserLoader) Load(userID string) (user User, shared bool, err error) {
	if l.fetch == nil {
		return User{}, false, fmt.Errorf("fetch user function is nil")
	}

	value, err, shared := l.group.Do(userID, func() (any, error) {
		return l.fetch(userID)
	})
	if err != nil {
		return User{}, shared, err
	}

	user, ok := value.(User)
	if !ok {
		return User{}, shared, fmt.Errorf("unexpected user result type %T", value)
	}

	return user, shared, nil
}
