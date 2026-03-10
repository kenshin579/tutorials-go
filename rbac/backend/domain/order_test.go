package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		name     string
		from     OrderStatus
		to       OrderStatus
		expected bool
	}{
		{"pending to confirmed", OrderStatusPending, OrderStatusConfirmed, true},
		{"pending to cancelled", OrderStatusPending, OrderStatusCancelled, true},
		{"confirmed to shipped", OrderStatusConfirmed, OrderStatusShipped, true},
		{"confirmed to cancelled", OrderStatusConfirmed, OrderStatusCancelled, true},
		{"shipped to completed", OrderStatusShipped, OrderStatusCompleted, true},
		{"pending to shipped (invalid)", OrderStatusPending, OrderStatusShipped, false},
		{"pending to completed (invalid)", OrderStatusPending, OrderStatusCompleted, false},
		{"completed to pending (invalid)", OrderStatusCompleted, OrderStatusPending, false},
		{"cancelled to pending (invalid)", OrderStatusCancelled, OrderStatusPending, false},
		{"shipped to cancelled (invalid)", OrderStatusShipped, OrderStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidTransition(tt.from, tt.to))
		})
	}
}

func TestCanTransition_Admin(t *testing.T) {
	tests := []struct {
		name     string
		from     OrderStatus
		to       OrderStatus
		expected bool
	}{
		{"pending to confirmed", OrderStatusPending, OrderStatusConfirmed, true},
		{"pending to cancelled", OrderStatusPending, OrderStatusCancelled, true},
		{"confirmed to shipped", OrderStatusConfirmed, OrderStatusShipped, true},
		{"confirmed to cancelled", OrderStatusConfirmed, OrderStatusCancelled, true},
		{"shipped to completed", OrderStatusShipped, OrderStatusCompleted, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CanTransition("admin", tt.from, tt.to))
		})
	}
}

func TestCanTransition_Manager(t *testing.T) {
	tests := []struct {
		name     string
		from     OrderStatus
		to       OrderStatus
		expected bool
	}{
		{"confirmed to shipped", OrderStatusConfirmed, OrderStatusShipped, true},
		{"pending to cancelled", OrderStatusPending, OrderStatusCancelled, true},
		{"confirmed to cancelled", OrderStatusConfirmed, OrderStatusCancelled, true},
		{"pending to confirmed (forbidden)", OrderStatusPending, OrderStatusConfirmed, false},
		{"shipped to completed (forbidden)", OrderStatusShipped, OrderStatusCompleted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CanTransition("manager", tt.from, tt.to))
		})
	}
}

func TestCanTransition_User(t *testing.T) {
	tests := []struct {
		name     string
		from     OrderStatus
		to       OrderStatus
		expected bool
	}{
		{"pending to cancelled", OrderStatusPending, OrderStatusCancelled, true},
		{"pending to confirmed (forbidden)", OrderStatusPending, OrderStatusConfirmed, false},
		{"confirmed to shipped (forbidden)", OrderStatusConfirmed, OrderStatusShipped, false},
		{"confirmed to cancelled (forbidden)", OrderStatusConfirmed, OrderStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CanTransition("user", tt.from, tt.to))
		})
	}
}

func TestCanTransitionWithRoles(t *testing.T) {
	// manager 단독으로는 pending→confirmed 불가, admin이 있으면 가능
	assert.False(t, CanTransitionWithRoles([]string{"manager"}, OrderStatusPending, OrderStatusConfirmed))
	assert.True(t, CanTransitionWithRoles([]string{"manager", "admin"}, OrderStatusPending, OrderStatusConfirmed))

	// user 단독으로는 confirmed→cancelled 불가
	assert.False(t, CanTransitionWithRoles([]string{"user"}, OrderStatusConfirmed, OrderStatusCancelled))

	// unknown role
	assert.False(t, CanTransitionWithRoles([]string{"unknown"}, OrderStatusPending, OrderStatusCancelled))
}
