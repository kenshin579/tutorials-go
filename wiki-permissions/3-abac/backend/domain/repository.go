package domain

// UserRepository는 User 엔티티에 대한 데이터 액세스 인터페이스다.
type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Create(u *User) error
}

// PageRepository는 Page 엔티티에 대한 데이터 액세스 인터페이스다.
// ABAC에서는 SQL 단계에서 권한 필터링을 하지 않고 List로 모두 가져온 뒤
// usecase 단에서 정책 평가로 필터링한다.
type PageRepository interface {
	FindByID(id uint) (*Page, error)
	List() ([]Page, error)
	Create(p *Page) error
	Update(p *Page) error
}

// DepartmentRepository는 Department 엔티티에 대한 데이터 액세스 인터페이스다.
type DepartmentRepository interface {
	FindByID(id uint) (*Department, error)
	List() ([]Department, error)
}

// ErrNotFound는 리소스(User/Page/Department 등)를 찾지 못했을 때 반환되는 에러다.
type ErrNotFound struct{ Resource string }

func (e ErrNotFound) Error() string { return e.Resource + " not found" }
