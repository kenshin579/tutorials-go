package domain

// UserRepository는 User 엔티티에 대한 데이터 액세스 인터페이스다.
type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Create(u *User) error
}

// PageRepository는 Page 엔티티에 대한 데이터 액세스 인터페이스다.
type PageRepository interface {
	FindByID(id uint) (*Page, error)
	ListAccessibleBy(userID uint) ([]Page, error)
	Update(p *Page) error
	Create(p *Page) error
}

// ACLRepository는 ACLEntry에 대한 데이터 액세스 인터페이스다.
type ACLRepository interface {
	FindByPageAndUser(pageID, userID uint) ([]ACLEntry, error)
	ListByPage(pageID uint) ([]ACLEntry, error)
	Grant(pageID, userID uint, action Action) error
	Revoke(pageID, userID uint) error
}

// ErrNotFound는 리소스(User/Page 등)를 찾지 못했을 때 반환되는 에러다.
type ErrNotFound struct{ Resource string }

func (e ErrNotFound) Error() string { return e.Resource + " not found" }
