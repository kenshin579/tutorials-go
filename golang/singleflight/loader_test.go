package singleflight

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const concurrentRequests = 10

type loadResult struct {
	user   User
	shared bool
	err    error
}

// TestWithoutSingleflightCallsDataSourceForEveryRequest는 중복 제거를 적용하지 않은
// 상태에서 같은 사용자를 동시에 요청하면 데이터 소스도 요청 수만큼 호출됨을 확인한다.
func TestWithoutSingleflightCallsDataSourceForEveryRequest(t *testing.T) {
	var calls atomic.Int64
	// started는 각 goroutine이 fetch에 진입했음을 테스트 본체에 알린다.
	// 모든 진입 신호를 막힘없이 담을 수 있도록 요청 수만큼 버퍼를 둔다.
	started := make(chan struct{}, concurrentRequests)
	// fetch는 release가 닫힐 때까지 대기하므로 느린 데이터 소스처럼 동작한다.
	release := make(chan struct{})

	fetch := func(userID string) (User, error) {
		calls.Add(1)
		started <- struct{}{} // fetch 진입을 알린다. 값 자체에는 의미가 없다.
		<-release             // close(release)가 호출될 때까지 조회 완료를 막는다.
		return User{ID: userID, Name: "Gopher"}, nil
	}

	var wg sync.WaitGroup
	for range concurrentRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = fetch("user-1")
		}()
	}

	// 모든 fetch가 동시에 실행 중인 상태가 될 때까지 진입 신호를 받는다.
	for range concurrentRequests {
		<-started
	}
	// 채널을 닫으면 release를 기다리는 모든 goroutine이 한꺼번에 깨어난다.
	close(release)
	wg.Wait()

	if got := calls.Load(); got != concurrentRequests {
		t.Fatalf("fetch call count = %d, want %d", got, concurrentRequests)
	}
}

// TestUserLoaderCombinesConcurrentCallsForSameKey는 같은 key의 동시 요청이 하나의
// 데이터 소스 호출로 합쳐지고 모든 호출자가 같은 결과를 공유하는지 확인한다.
func TestUserLoaderCombinesConcurrentCallsForSameKey(t *testing.T) {
	var calls atomic.Int64
	// singleflight가 제대로 동작하면 started에는 신호가 하나만 들어온다.
	started := make(chan struct{}, concurrentRequests)
	// 최초 fetch를 멈춰 둔 사이 나머지 요청들이 같은 호출에 합류하게 한다.
	release := make(chan struct{})

	loader := NewUserLoader(func(userID string) (User, error) {
		calls.Add(1)
		started <- struct{}{} // 실제 데이터 소스 호출이 시작됐음을 알린다.
		<-release
		return User{ID: userID, Name: "Gopher"}, nil
	})

	// 반환 채널을 먼저 받고, 모든 요청 결과는 아래에서 차례로 수집한다.
	results := startConcurrentLoads(loader, concurrentRequests, "user-1")
	<-started // 최초 요청이 fetch 안에서 대기 중임을 확인한다.
	waitForDuplicates(t, results)
	// 최초 fetch를 완료시키면 같은 key로 대기하던 요청들도 같은 결과를 받는다.
	close(release)

	for range concurrentRequests {
		result := <-results
		if result.err != nil {
			t.Fatalf("Load() error = %v", result.err)
		}
		if result.user.ID != "user-1" {
			t.Errorf("Load() user ID = %q, want %q", result.user.ID, "user-1")
		}
		if !result.shared {
			t.Error("Load() shared = false, want true")
		}
	}

	if got := calls.Load(); got != 1 {
		t.Fatalf("fetch call count = %d, want 1", got)
	}
}

// TestUserLoaderRunsDifferentKeysIndependently는 서로 다른 key의 요청이 하나로
// 합쳐지지 않고 각각 독립적으로 데이터 소스를 호출하는지 확인한다.
func TestUserLoaderRunsDifferentKeysIndependently(t *testing.T) {
	var calls atomic.Int64
	// 어떤 key의 fetch가 시작됐는지 확인하기 위해 사용자 ID를 전달한다.
	started := make(chan string, 2)
	// 두 key의 fetch가 모두 시작된 것을 확인한 뒤 함께 완료시킨다.
	release := make(chan struct{})

	loader := NewUserLoader(func(userID string) (User, error) {
		calls.Add(1)
		started <- userID
		<-release
		return User{ID: userID}, nil
	})

	// key별 결과 채널을 분리해 각 요청의 반환값을 확인한다.
	result1 := make(chan loadResult, 1)
	result2 := make(chan loadResult, 1)
	go loadInto(loader, "user-1", result1)
	go loadInto(loader, "user-2", result2)

	// started에서 두 ID를 받았다는 것은 서로 다른 두 fetch가 실행됐다는 뜻이다.
	seen := map[string]bool{<-started: true, <-started: true}
	close(release)

	for _, result := range []loadResult{<-result1, <-result2} {
		if result.err != nil {
			t.Fatalf("Load() error = %v", result.err)
		}
		if result.shared {
			t.Error("Load() shared = true for a key requested only once")
		}
	}

	if !seen["user-1"] || !seen["user-2"] {
		t.Fatalf("started keys = %v, want user-1 and user-2", seen)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("fetch call count = %d, want 2", got)
	}
}

// TestUserLoaderDoesNotCacheCompletedResult는 첫 호출이 끝난 뒤 같은 key를 다시
// 요청하면 데이터 소스가 재실행되어 singleflight가 캐시가 아님을 확인한다.
func TestUserLoaderDoesNotCacheCompletedResult(t *testing.T) {
	var calls atomic.Int64
	loader := NewUserLoader(func(userID string) (User, error) {
		calls.Add(1)
		return User{ID: userID}, nil
	})

	for range 2 {
		_, shared, err := loader.Load("user-1")
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if shared {
			t.Error("sequential Load() shared = true, want false")
		}
	}

	if got := calls.Load(); got != 2 {
		t.Fatalf("fetch call count = %d, want 2", got)
	}
}

// TestUserLoaderSharesErrorWithConcurrentCallers는 데이터 소스에서 발생한 오류도
// 같은 key를 기다리는 모든 동시 호출자에게 공유되는지 확인한다.
func TestUserLoaderSharesErrorWithConcurrentCallers(t *testing.T) {
	var calls atomic.Int64
	// 정상 결과 테스트와 마찬가지로 최초 fetch를 멈춰 중복 요청이 합류하게 한다.
	started := make(chan struct{}, concurrentRequests)
	release := make(chan struct{})
	wantErr := errors.New("data source unavailable")

	loader := NewUserLoader(func(string) (User, error) {
		calls.Add(1)
		started <- struct{}{}
		<-release
		return User{}, wantErr
	})

	results := startConcurrentLoads(loader, concurrentRequests, "user-1")
	<-started
	waitForDuplicates(t, results)
	close(release)

	for range concurrentRequests {
		result := <-results
		if !errors.Is(result.err, wantErr) {
			t.Errorf("Load() error = %v, want %v", result.err, wantErr)
		}
		if !result.shared {
			t.Error("Load() shared = false, want true")
		}
	}

	if got := calls.Load(); got != 1 {
		t.Fatalf("fetch call count = %d, want 1", got)
	}
}

func startConcurrentLoads(loader *UserLoader, count int, userID string) <-chan loadResult {
	// start는 모든 goroutine을 같은 시점에 출발시키는 시작 신호다.
	// 각 goroutine은 채널이 닫히기 전까지 아래의 수신 연산에서 대기한다.
	start := make(chan struct{})
	// 호출자가 결과를 즉시 읽지 않아도 goroutine이 송신에서 막히지 않도록
	// 요청 수만큼 버퍼를 둔다.
	results := make(chan loadResult, count)

	for range count {
		go func() {
			<-start // close(start)가 모든 goroutine을 동시에 깨운다.
			loadInto(loader, userID, results)
		}()
	}

	// 값을 여러 번 보내는 대신 채널을 한 번 닫아 모든 대기자를 출발시킨다.
	close(start)
	return results
}

func loadInto(loader *UserLoader, userID string, results chan<- loadResult) {
	user, shared, err := loader.Load(userID)
	// 송신 전용 채널을 사용해 이 함수가 results를 수신하지 않음을 표현한다.
	results <- loadResult{user: user, shared: shared, err: err}
}

// The leader remains blocked while duplicate goroutines enter Group.Do. The
// timeout also makes a scheduler failure explicit instead of hanging the test.
func waitForDuplicates(t *testing.T, results <-chan loadResult) {
	t.Helper()

	select {
	// release가 닫히기 전에 결과가 오면 fetch가 예상과 달리 먼저 끝난 것이다.
	case result := <-results:
		t.Fatalf("Load() returned before data source was released: %+v", result)
	// 짧게 기다리는 동안 중복 goroutine들이 Group.Do에 진입할 시간을 준다.
	case <-time.After(20 * time.Millisecond):
	}
}
