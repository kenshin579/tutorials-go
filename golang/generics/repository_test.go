package go_generics

import (
	"errors"
	"fmt"
	"sync"
)

// ============================================================
// Generic Repository 패턴
// ============================================================

// Identifiable - 모든 엔티티가 구현해야 하는 인터페이스
type Identifiable interface {
	GetID() string
}

// MemoryStore - Generic in-memory repository
type MemoryStore[T Identifiable] struct {
	mu   sync.RWMutex
	data map[string]T
}

// NewMemoryStore - 새 MemoryStore 생성
func NewMemoryStore[T Identifiable]() *MemoryStore[T] {
	return &MemoryStore[T]{
		data: make(map[string]T),
	}
}

func (s *MemoryStore[T]) Save(entity T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[entity.GetID()] = entity
}

func (s *MemoryStore[T]) FindByID(id string) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entity, ok := s.data[id]
	if !ok {
		var zero T
		return zero, errors.New("not found")
	}
	return entity, nil
}

func (s *MemoryStore[T]) FindAll() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, 0, len(s.data))
	for _, v := range s.data {
		result = append(result, v)
	}
	return result
}

func (s *MemoryStore[T]) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return false
	}
	delete(s.data, id)
	return true
}

func (s *MemoryStore[T]) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// ============================================================
// 엔티티 정의
// ============================================================

// UserEntity - 사용자 엔티티
type UserEntity struct {
	ID    string
	Name  string
	Email string
}

func (u UserEntity) GetID() string { return u.ID }

// ProductEntity - 상품 엔티티
type ProductEntity struct {
	ID    string
	Name  string
	Price int
}

func (p ProductEntity) GetID() string { return p.ID }

// ============================================================
// Example: 동일한 Repository를 다른 타입에 재사용
// ============================================================

func Example_genericRepository() {
	// 사용자 저장소
	userStore := NewMemoryStore[UserEntity]()
	userStore.Save(UserEntity{ID: "u1", Name: "Alice", Email: "alice@example.com"})
	userStore.Save(UserEntity{ID: "u2", Name: "Bob", Email: "bob@example.com"})

	user, _ := userStore.FindByID("u1")
	fmt.Println(user.Name, user.Email)
	fmt.Println("users:", userStore.Count())

	// 상품 저장소 - 동일한 MemoryStore를 재사용
	productStore := NewMemoryStore[ProductEntity]()
	productStore.Save(ProductEntity{ID: "p1", Name: "Laptop", Price: 1500})
	productStore.Save(ProductEntity{ID: "p2", Name: "Mouse", Price: 30})

	product, _ := productStore.FindByID("p2")
	fmt.Println(product.Name, product.Price)
	fmt.Println("products:", productStore.Count())

	// 삭제
	productStore.Delete("p1")
	fmt.Println("after delete:", productStore.Count())

	// Output:
	// Alice alice@example.com
	// users: 2
	// Mouse 30
	// products: 2
	// after delete: 1
}

// ============================================================
// 타입 안전 컬렉션: Generic Set
// ============================================================

// Set - 타입 안전한 집합
type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

func (s *Set[T]) Has(item T) bool {
	_, ok := s.items[item]
	return ok
}

func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func Example_genericSet() {
	// string Set
	tags := NewSet[string]()
	tags.Add("go")
	tags.Add("generics")
	tags.Add("go") // 중복 무시
	fmt.Println(tags.Has("go"))
	fmt.Println(tags.Has("rust"))
	fmt.Println(tags.Size())

	// int Set
	ids := NewSet[int]()
	ids.Add(1)
	ids.Add(2)
	ids.Add(3)
	ids.Remove(2)
	fmt.Println(ids.Size())

	// Output:
	// true
	// false
	// 2
	// 2
}
