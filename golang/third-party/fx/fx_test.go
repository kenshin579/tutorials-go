package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// --- 테스트용 인터페이스 및 구현체 ---

type Logger interface {
	Log(msg string)
}

type stdLogger struct{}

func (l *stdLogger) Log(msg string) { fmt.Println(msg) }

func NewLogger() Logger { return &stdLogger{} }

type UserRepository interface {
	FindByID(id int) string
}

type mysqlUserRepo struct {
	logger Logger
}

func (r *mysqlUserRepo) FindByID(id int) string {
	return fmt.Sprintf("user-%d", id)
}

func NewMysqlUserRepo(logger Logger) UserRepository {
	return &mysqlUserRepo{logger: logger}
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

type OrderRepository interface {
	FindByID(id int) string
}

type mysqlOrderRepo struct{}

func (r *mysqlOrderRepo) FindByID(id int) string {
	return fmt.Sprintf("order-%d", id)
}

func NewMysqlOrderRepo() OrderRepository {
	return &mysqlOrderRepo{}
}

type OrderService struct {
	repo   OrderRepository
	logger Logger
}

func NewOrderService(repo OrderRepository, logger Logger) *OrderService {
	return &OrderService{repo: repo, logger: logger}
}

// --- 기본 fx 패턴 ---

func TestFx_Provide_Invoke(t *testing.T) {
	var svc *UserService

	app := fxtest.New(t,
		fx.Provide(
			NewLogger,
			NewMysqlUserRepo,
			NewUserService,
		),
		fx.Invoke(func(s *UserService) {
			svc = s
		}),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.NotNil(t, svc)
	assert.Equal(t, "user-1", svc.repo.FindByID(1))
}

func TestFx_Supply(t *testing.T) {
	// fx.Supply: 이미 생성된 값을 직접 제공
	type Config struct {
		DBHost string
		DBPort int
	}

	cfg := &Config{DBHost: "localhost", DBPort: 3306}

	var received *Config
	app := fxtest.New(t,
		fx.Supply(cfg),
		fx.Invoke(func(c *Config) {
			received = c
		}),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Equal(t, "localhost", received.DBHost)
	assert.Equal(t, 3306, received.DBPort)
}

// --- fx.Module 패턴 ---

func TestFx_Module(t *testing.T) {
	// 도메인별 의존성을 Module로 그룹화
	var UserModule = fx.Module("user",
		fx.Provide(
			NewMysqlUserRepo,
			NewUserService,
		),
	)

	var OrderModule = fx.Module("order",
		fx.Provide(
			NewMysqlOrderRepo,
			NewOrderService,
		),
	)

	var userSvc *UserService
	var orderSvc *OrderService

	app := fxtest.New(t,
		fx.Provide(NewLogger), // 공통 의존성
		UserModule,
		OrderModule,
		fx.Invoke(func(u *UserService, o *OrderService) {
			userSvc = u
			orderSvc = o
		}),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Equal(t, "user-1", userSvc.repo.FindByID(1))
	assert.Equal(t, "order-42", orderSvc.repo.FindByID(42))
}

// --- fx.Decorate 패턴 ---

// 로깅 데코레이터: 기존 UserRepository를 래핑하여 호출 로그를 추가
type loggingUserRepo struct {
	inner  UserRepository
	logger Logger
	calls  []string
}

func (r *loggingUserRepo) FindByID(id int) string {
	r.calls = append(r.calls, fmt.Sprintf("FindByID(%d)", id))
	return r.inner.FindByID(id)
}

func TestFx_Decorate(t *testing.T) {
	decorated := &loggingUserRepo{}

	app := fxtest.New(t,
		fx.Provide(
			NewLogger,
			NewMysqlUserRepo,
			NewUserService,
		),
		// Decorate: 기존 UserRepository를 로깅 래퍼로 교체
		fx.Decorate(func(repo UserRepository, logger Logger) UserRepository {
			decorated.inner = repo
			decorated.logger = logger
			return decorated
		}),
		fx.Invoke(func(svc *UserService) {
			svc.repo.FindByID(1)
			svc.repo.FindByID(2)
		}),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Len(t, decorated.calls, 2)
	assert.Equal(t, "FindByID(1)", decorated.calls[0])
	assert.Equal(t, "FindByID(2)", decorated.calls[1])
}

// --- fx.Annotate + Named 의존성 ---

type DBConnection struct {
	Name string
	DSN  string
}

func NewReadDB() *DBConnection {
	return &DBConnection{Name: "read", DSN: "read-replica:3306"}
}

func NewWriteDB() *DBConnection {
	return &DBConnection{Name: "write", DSN: "primary:3306"}
}

// fx.In을 사용한 Named 의존성 주입
type DBParams struct {
	fx.In
	ReadDB  *DBConnection `name:"readDB"`
	WriteDB *DBConnection `name:"writeDB"`
}

type DBService struct {
	readDB  *DBConnection
	writeDB *DBConnection
}

func NewDBService(params DBParams) *DBService {
	return &DBService{
		readDB:  params.ReadDB,
		writeDB: params.WriteDB,
	}
}

func TestFx_Annotate_Named(t *testing.T) {
	var svc *DBService

	app := fxtest.New(t,
		fx.Provide(
			// fx.Annotate로 동일 타입에 이름 부여
			fx.Annotate(NewReadDB, fx.ResultTags(`name:"readDB"`)),
			fx.Annotate(NewWriteDB, fx.ResultTags(`name:"writeDB"`)),
			NewDBService,
		),
		fx.Invoke(func(s *DBService) {
			svc = s
		}),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Equal(t, "read-replica:3306", svc.readDB.DSN)
	assert.Equal(t, "primary:3306", svc.writeDB.DSN)
}

// --- fx.Replace (테스트 시 Mock 주입) ---

type mockUserRepo struct{}

func (r *mockUserRepo) FindByID(id int) string {
	return fmt.Sprintf("mock-user-%d", id)
}

func TestFx_Replace_Mock(t *testing.T) {
	var svc *UserService

	app := fxtest.New(t,
		fx.Provide(
			NewLogger,
			NewMysqlUserRepo,
			NewUserService,
		),
		// Replace: 테스트 시 실제 구현을 Mock으로 교체
		fx.Replace(fx.Annotate(&mockUserRepo{}, fx.As(new(UserRepository)))),
		fx.Invoke(func(s *UserService) {
			svc = s
		}),
	)
	defer app.RequireStop()
	app.RequireStart()

	// Mock이 주입되었으므로 "mock-user-" 접두사
	assert.Equal(t, "mock-user-1", svc.repo.FindByID(1))
}

// --- fx.Lifecycle ---

func TestFx_Lifecycle(t *testing.T) {
	var startCalled, stopCalled bool

	app := fxtest.New(t,
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					startCalled = true
					return nil
				},
				OnStop: func(context.Context) error {
					stopCalled = true
					return nil
				},
			})
		}),
	)

	app.RequireStart()
	assert.True(t, startCalled)
	assert.False(t, stopCalled)

	app.RequireStop()
	assert.True(t, stopCalled)
}

// --- fx.Group: 동일 인터페이스 여러 구현체 모으기 ---

type Notifier interface {
	Send(msg string) string
}

type EmailNotifier struct{}

func (e *EmailNotifier) Send(msg string) string { return "email:" + msg }

type SlackNotifier struct{}

func (s *SlackNotifier) Send(msg string) string { return "slack:" + msg }

type SMSNotifier struct{}

func (s *SMSNotifier) Send(msg string) string { return "sms:" + msg }

// fx.In의 group 태그로 같은 그룹의 모든 구현체를 슬라이스로 수신
type NotifierParams struct {
	fx.In
	Notifiers []Notifier `group:"notifiers"`
}

type NotifierService struct {
	notifiers []Notifier
}

func NewNotifierService(p NotifierParams) *NotifierService {
	return &NotifierService{notifiers: p.Notifiers}
}

func TestFx_Group_ValueGroups(t *testing.T) {
	var svc *NotifierService

	app := fxtest.New(t,
		fx.Provide(
			// fx.ResultTags로 같은 group에 여러 구현체 등록
			fx.Annotate(func() Notifier { return &EmailNotifier{} },
				fx.ResultTags(`group:"notifiers"`)),
			fx.Annotate(func() Notifier { return &SlackNotifier{} },
				fx.ResultTags(`group:"notifiers"`)),
			fx.Annotate(func() Notifier { return &SMSNotifier{} },
				fx.ResultTags(`group:"notifiers"`)),
			NewNotifierService,
		),
		fx.Invoke(func(s *NotifierService) {
			svc = s
		}),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Len(t, svc.notifiers, 3)

	var results []string
	for _, n := range svc.notifiers {
		results = append(results, n.Send("hi"))
	}
	assert.Contains(t, results, "email:hi")
	assert.Contains(t, results, "slack:hi")
	assert.Contains(t, results, "sms:hi")
}

// --- fx.Populate: fx 컨테이너에서 인스턴스를 외부 변수로 추출 ---

func TestFx_Populate(t *testing.T) {
	// 형태 1: 단일 인스턴스 추출
	var svc1 *UserService
	app1 := fxtest.New(t,
		fx.Provide(NewLogger, NewMysqlUserRepo, NewUserService),
		fx.Populate(&svc1),
	)
	defer app1.RequireStop()
	app1.RequireStart()

	assert.NotNil(t, svc1)
	assert.Equal(t, "user-1", svc1.repo.FindByID(1))

	// 형태 2: 여러 인스턴스를 한꺼번에 추출
	var (
		svc2    *UserService
		logger2 Logger
	)
	app2 := fxtest.New(t,
		fx.Provide(NewLogger, NewMysqlUserRepo, NewUserService),
		fx.Populate(&svc2, &logger2),
	)
	defer app2.RequireStop()
	app2.RequireStart()

	assert.NotNil(t, svc2)
	assert.NotNil(t, logger2)
}

// --- fx.Private: Module 내부 전용 의존성 캡슐화 ---

type internalDB struct {
	name string
}

func newInternalDB() *internalDB {
	return &internalDB{name: "private-db"}
}

type ModuleService struct {
	db *internalDB
}

func newModuleService(db *internalDB) *ModuleService {
	return &ModuleService{db: db}
}

func TestFx_Private(t *testing.T) {
	// PrivateModule:
	//   - 첫 번째 fx.Provide()에 fx.Private을 마지막 인자로 넣어 *internalDB를 Module 내부 전용으로
	//   - 두 번째 fx.Provide()는 일반 노출. *ModuleService는 외부에서 추출 가능
	PrivateModule := fx.Module("private",
		fx.Provide(
			newInternalDB,
			fx.Private,
		),
		fx.Provide(newModuleService),
	)

	// 정상: ModuleService는 외부 노출 → 추출 성공
	var svc *ModuleService
	app := fxtest.New(t,
		PrivateModule,
		fx.Populate(&svc),
	)
	defer app.RequireStop()
	app.RequireStart()
	assert.Equal(t, "private-db", svc.db.name)

	// 비정상: *internalDB는 Module 내부 전용 → 외부 추출 시도 시 fx.New가 에러 반환
	var leaked *internalDB
	leakApp := fx.New(
		PrivateModule,
		fx.Populate(&leaked),
		fx.NopLogger, // 에러를 stdout으로 출력하지 않음
	)
	err := leakApp.Err()
	assert.Error(t, err, "internalDB는 Module 외부에서 보이지 않아야 한다")
	assert.Contains(t, err.Error(), "*main.internalDB",
		"에러 메시지에 외부 추출이 막힌 타입이 명시되어야 한다")
}

// --- fx.As: 구체 타입 생성자를 인터페이스로 등록 ---

type Cache interface {
	Get(key string) string
}

type Closer interface {
	Close() error
}

// RedisCache는 Cache와 Closer를 모두 구현한다.
type RedisCache struct {
	closed bool
}

func (c *RedisCache) Get(key string) string { return "redis:" + key }

func (c *RedisCache) Close() error {
	c.closed = true
	return nil
}

// NewRedisCache는 인터페이스가 아니라 구체 타입(*RedisCache)을 반환한다.
func NewRedisCache() *RedisCache { return &RedisCache{} }

type CacheService struct {
	cache Cache
}

// NewCacheService는 구체 타입이 아니라 Cache 인터페이스를 요구한다.
func NewCacheService(cache Cache) *CacheService {
	return &CacheService{cache: cache}
}

func TestFx_As_WithoutAnnotation(t *testing.T) {
	// fx는 타입으로 매칭한다. NewRedisCache를 그대로 Provide하면 그래프에는
	// *RedisCache만 등록되므로, Cache를 요구하는 NewCacheService는 짝을 못 찾는다.
	var svc *CacheService
	app := fx.New(
		fx.Provide(NewRedisCache, NewCacheService),
		fx.Populate(&svc),
		fx.NopLogger,
	)

	err := app.Err()
	assert.Error(t, err, "Cache 타입을 제공하는 생성자가 없어야 한다")
	assert.Contains(t, err.Error(), "main.Cache",
		"에러 메시지에 채우지 못한 인터페이스 타입이 명시되어야 한다")
}

func TestFx_As_ProvideAsInterface(t *testing.T) {
	// fx.As(new(Cache)): *RedisCache를 Cache 인터페이스로 등록
	var svc *CacheService

	app := fxtest.New(t,
		fx.Provide(
			fx.Annotate(NewRedisCache, fx.As(new(Cache))),
			NewCacheService,
		),
		fx.Populate(&svc),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Equal(t, "redis:user:1", svc.cache.Get("user:1"))
}

func TestFx_As_ReplacesOriginalType(t *testing.T) {
	// fx.As는 원래 반환 타입을 "추가"하는 게 아니라 "대체"한다.
	// Cache로 등록된 순간 *RedisCache는 그래프에서 사라진다.
	var concrete *RedisCache
	app := fx.New(
		fx.Provide(fx.Annotate(NewRedisCache, fx.As(new(Cache)))),
		fx.Populate(&concrete),
		fx.NopLogger,
	)

	err := app.Err()
	assert.Error(t, err, "fx.As로 대체된 구체 타입은 주입받을 수 없어야 한다")
	assert.Contains(t, err.Error(), "*main.RedisCache")
}

func TestFx_As_Self(t *testing.T) {
	// fx.As(fx.Self())를 함께 붙이면 원래 구체 타입도 그래프에 남는다 (v1.20+)
	var (
		cache    Cache
		concrete *RedisCache
	)

	app := fxtest.New(t,
		fx.Provide(fx.Annotate(NewRedisCache,
			fx.As(new(Cache)),
			fx.As(fx.Self()),
		)),
		fx.Populate(&cache, &concrete),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Equal(t, "redis:k", cache.Get("k"))
	assert.Same(t, concrete, cache, "두 타입은 같은 인스턴스를 가리켜야 한다")
}

func TestFx_As_MultipleInterfaces(t *testing.T) {
	// fx.As를 여러 번 붙이면 하나의 생성자를 여러 인터페이스로 등록한다
	var (
		cache  Cache
		closer Closer
	)

	app := fxtest.New(t,
		fx.Provide(fx.Annotate(NewRedisCache,
			fx.As(new(Cache)),
			fx.As(new(Closer)),
		)),
		fx.Populate(&cache, &closer),
	)
	defer app.RequireStop()
	app.RequireStart()

	assert.Equal(t, "redis:k", cache.Get("k"))
	assert.NoError(t, closer.Close())
	assert.Same(t, cache, closer, "두 인터페이스는 같은 인스턴스를 가리켜야 한다")
	assert.True(t, cache.(*RedisCache).closed, "Close()가 같은 인스턴스에 반영되어야 한다")
}
