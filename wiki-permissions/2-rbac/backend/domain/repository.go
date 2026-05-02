package domain

// UserRepository는 User 엔티티에 대한 데이터 액세스 인터페이스다.
type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	List() ([]User, error)
	Create(u *User) error
	AssignRole(userID, roleID uint) error
	RevokeRole(userID, roleID uint) error
}

// PageRepository는 Page 엔티티에 대한 데이터 액세스 인터페이스다.
type PageRepository interface {
	FindByID(id uint) (*Page, error)
	List() ([]Page, error)
	Create(p *Page) error
	Update(p *Page) error
	Delete(id uint) error
}

// RoleRepository는 Role 엔티티에 대한 데이터 액세스 인터페이스다.
type RoleRepository interface {
	FindByID(id uint) (*Role, error)
	FindByName(name string) (*Role, error)
	List() ([]Role, error)
}

// PermissionRepository는 Permission에 대한 데이터 액세스 인터페이스다.
type PermissionRepository interface {
	// FindByUserID는 사용자 → role → permission 3-hop JOIN 결과를 중복 제거하여 반환한다.
	FindByUserID(userID uint) ([]Permission, error)
}

// ErrNotFound는 리소스(User/Page/Role 등)를 찾지 못했을 때 반환되는 에러다.
type ErrNotFound struct{ Resource string }

func (e ErrNotFound) Error() string { return e.Resource + " not found" }
