package go_generics

import (
	"fmt"
	"sort"
)

// ============================================================
// 1. Before: interface{} 기반 Config
// ============================================================

// ConfigOld - interface{} 기반 설정 저장소 (타입 안전하지 않음)
type ConfigOld struct {
	data map[string]interface{}
}

func NewConfigOld() *ConfigOld {
	return &ConfigOld{data: make(map[string]interface{})}
}

func (c *ConfigOld) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *ConfigOld) Get(key string) (interface{}, bool) {
	v, ok := c.data[key]
	return v, ok
}

func Example_configOld() {
	cfg := NewConfigOld()
	cfg.Set("port", 8080)
	cfg.Set("host", "localhost")

	// 사용할 때 타입 단언 필요 → 런타임 에러 위험
	if v, ok := cfg.Get("port"); ok {
		port := v.(int) // 타입 단언 필요
		fmt.Println("port:", port)
	}
	if v, ok := cfg.Get("host"); ok {
		host := v.(string) // 타입 단언 필요
		fmt.Println("host:", host)
	}

	// Output:
	// port: 8080
	// host: localhost
}

// ============================================================
// 2. After: Generics 기반 Config
// ============================================================

// TypedConfig - 타입 안전한 설정 항목
type TypedConfig[T any] struct {
	key          string
	defaultValue T
}

// NewTypedConfig - 새 설정 항목 생성
func NewTypedConfig[T any](key string, defaultValue T) *TypedConfig[T] {
	return &TypedConfig[T]{key: key, defaultValue: defaultValue}
}

func (tc *TypedConfig[T]) Key() string     { return tc.key }
func (tc *TypedConfig[T]) Default() T      { return tc.defaultValue }

// ConfigStore - 설정 저장소
type ConfigStore struct {
	data map[string]interface{} // 내부 저장은 interface{} 사용
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{data: make(map[string]interface{})}
}

// SetConfig - 타입 안전하게 설정 값 저장
func SetConfig[T any](store *ConfigStore, cfg *TypedConfig[T], value T) {
	store.data[cfg.Key()] = value
}

// GetConfig - 타입 안전하게 설정 값 조회 (타입 단언 불필요)
func GetConfig[T any](store *ConfigStore, cfg *TypedConfig[T]) T {
	if v, ok := store.data[cfg.Key()]; ok {
		return v.(T)
	}
	return cfg.Default()
}

func Example_configNew() {
	// 설정 항목을 타입과 함께 정의
	portCfg := NewTypedConfig("port", 3000)
	hostCfg := NewTypedConfig("host", "0.0.0.0")
	debugCfg := NewTypedConfig("debug", false)

	store := NewConfigStore()
	SetConfig(store, portCfg, 8080)
	SetConfig(store, hostCfg, "localhost")

	// 타입 단언 불필요 - 컴파일 타임에 타입 보장
	port := GetConfig(store, portCfg)
	host := GetConfig(store, hostCfg)
	debug := GetConfig(store, debugCfg) // 설정 없으면 기본값 반환

	fmt.Println("port:", port)
	fmt.Println("host:", host)
	fmt.Println("debug:", debug)

	// Output:
	// port: 8080
	// host: localhost
	// debug: false
}

// ============================================================
// 3. Before: interface{} 기반 이벤트 버스
// ============================================================

// EventBusOld - interface{} 기반 이벤트 버스
type EventBusOld struct {
	handlers map[string][]func(interface{})
}

func NewEventBusOld() *EventBusOld {
	return &EventBusOld{handlers: make(map[string][]func(interface{}))}
}

func (bus *EventBusOld) On(event string, handler func(interface{})) {
	bus.handlers[event] = append(bus.handlers[event], handler)
}

func (bus *EventBusOld) Emit(event string, data interface{}) {
	for _, handler := range bus.handlers[event] {
		handler(data)
	}
}

func Example_eventBusOld() {
	bus := NewEventBusOld()

	bus.On("user.created", func(data interface{}) {
		// 타입 단언 필요 → 런타임 에러 위험
		name := data.(string)
		fmt.Println("created:", name)
	})

	bus.Emit("user.created", "Alice")

	// Output:
	// created: Alice
}

// ============================================================
// 4. After: Generics 기반 타입 안전 이벤트 핸들러
// ============================================================

// TypedHandler - 타입 안전한 이벤트 핸들러
type TypedHandler[T any] struct {
	handlers []func(T)
}

func NewTypedHandler[T any]() *TypedHandler[T] {
	return &TypedHandler[T]{}
}

func (h *TypedHandler[T]) On(handler func(T)) {
	h.handlers = append(h.handlers, handler)
}

func (h *TypedHandler[T]) Emit(data T) {
	for _, handler := range h.handlers {
		handler(data)
	}
}

// UserCreatedEvent - 구체적인 이벤트 타입
type UserCreatedEvent struct {
	Name  string
	Email string
}

func Example_eventHandlerNew() {
	handler := NewTypedHandler[UserCreatedEvent]()

	handler.On(func(e UserCreatedEvent) {
		// 타입 단언 불필요 - 컴파일 타임에 타입 보장
		fmt.Printf("created: %s (%s)\n", e.Name, e.Email)
	})

	handler.Emit(UserCreatedEvent{Name: "Alice", Email: "alice@example.com"})

	// 컴파일 에러: 잘못된 타입은 컴파일 시 차단
	// handler.Emit("wrong type") → 컴파일 에러

	// Output:
	// created: Alice (alice@example.com)
}

// ============================================================
// 5. Before → After: 정렬 함수 마이그레이션
// ============================================================

// sortInterfaceSlice - interface{} 기반 (Before)
func sortInterfaceSlice(items []interface{}, less func(a, b interface{}) bool) {
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})
}

// SortSliceBy - Generics 기반 (After)
func SortSliceBy[T any](items []T, less func(a, b T) bool) {
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})
}

func Example_sortMigration() {
	// Before: interface{} 기반 - 타입 단언 필요
	oldItems := []interface{}{3, 1, 2}
	sortInterfaceSlice(oldItems, func(a, b interface{}) bool {
		return a.(int) < b.(int) // 타입 단언 필요
	})
	fmt.Println(oldItems)

	// After: Generics 기반 - 타입 안전
	newItems := []int{3, 1, 2}
	SortSliceBy(newItems, func(a, b int) bool {
		return a < b // 타입 단언 불필요
	})
	fmt.Println(newItems)

	// Output:
	// [1 2 3]
	// [1 2 3]
}
