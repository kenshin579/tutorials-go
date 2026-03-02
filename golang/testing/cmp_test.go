package go_testing

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Address struct {
	City    string
	Country string
}

type Profile struct {
	User    User
	Address Address
	Tags    []string
}

// TestCmp_Equal - cmp.Equal로 구조체 비교
func TestCmp_Equal(t *testing.T) {
	user1 := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	user2 := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	if !cmp.Equal(user1, user2) {
		t.Errorf("users should be equal")
	}
}

// TestCmp_Diff - cmp.Diff로 차이점 확인
// 테스트 실패 시 어떤 필드가 다른지 읽기 쉽게 출력된다.
func TestCmp_Diff(t *testing.T) {
	want := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	got := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("User mismatch (-want +got):\n%s", diff)
	}
}

// TestCmp_IgnoreFields - cmpopts.IgnoreFields로 특정 필드 제외 비교
// CreatedAt, UpdatedAt 같은 타임스탬프 필드를 무시할 때 유용하다.
func TestCmp_IgnoreFields(t *testing.T) {
	now := time.Now()
	user1 := User{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: now}
	user2 := User{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: now.Add(time.Hour)}

	opts := cmpopts.IgnoreFields(User{}, "CreatedAt", "UpdatedAt")
	if diff := cmp.Diff(user1, user2, opts); diff != "" {
		t.Errorf("User mismatch (-want +got):\n%s", diff)
	}
}

// TestCmp_SortSlices - cmpopts.SortSlices로 슬라이스 순서 무관 비교
func TestCmp_SortSlices(t *testing.T) {
	want := []string{"banana", "apple", "cherry"}
	got := []string{"cherry", "banana", "apple"}

	opts := cmpopts.SortSlices(func(a, b string) bool { return a < b })
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("slices mismatch (-want +got):\n%s", diff)
	}
}

// TestCmp_EquateEmpty - cmpopts.EquateEmpty로 nil 슬라이스와 빈 슬라이스를 동일 취급
func TestCmp_EquateEmpty(t *testing.T) {
	var nilSlice []string
	emptySlice := []string{}

	opts := cmpopts.EquateEmpty()
	if diff := cmp.Diff(nilSlice, emptySlice, opts); diff != "" {
		t.Errorf("slices should be equal (-want +got):\n%s", diff)
	}
}

// TestCmp_NestedStruct - 중첩 구조체 비교
func TestCmp_NestedStruct(t *testing.T) {
	want := Profile{
		User:    User{ID: 1, Name: "Alice", Email: "alice@example.com"},
		Address: Address{City: "Seoul", Country: "KR"},
		Tags:    []string{"go", "testing"},
	}
	got := Profile{
		User:    User{ID: 1, Name: "Alice", Email: "alice@example.com"},
		Address: Address{City: "Seoul", Country: "KR"},
		Tags:    []string{"testing", "go"},
	}

	opts := cmpopts.SortSlices(func(a, b string) bool { return a < b })
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Profile mismatch (-want +got):\n%s", diff)
	}
}

// TestCmp_StructSlice - 구조체 슬라이스 비교 (ID 기준 정렬)
func TestCmp_StructSlice(t *testing.T) {
	want := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	got := []User{
		{ID: 3, Name: "Charlie"},
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	opts := cmpopts.SortSlices(func(a, b User) bool { return a.ID < b.ID })
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("users mismatch (-want +got):\n%s", diff)
	}
}
