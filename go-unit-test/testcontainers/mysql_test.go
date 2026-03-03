package go_testcontainers

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

type mysqlTestContainerTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(mysqlTestContainerTestSuite))
}

func (s *mysqlTestContainerTestSuite) SetupSuite() {
	s.db = NewMysqlDB()
	s.NoError(s.db.AutoMigrate(&User{}))
}

func (s *mysqlTestContainerTestSuite) TearDownTest() {
	s.db.Exec("DELETE FROM users")
}

func (s *mysqlTestContainerTestSuite) Test_CreateAndFind() {
	user := User{Name: "Frank", Email: "frank@example.com"}
	result := s.db.Create(&user)
	s.NoError(result.Error)
	s.NotZero(user.ID)

	var found User
	s.NoError(s.db.First(&found, user.ID).Error)
	s.Equal("Frank", found.Name)
	s.Equal("frank@example.com", found.Email)
}

func (s *mysqlTestContainerTestSuite) Test_Update() {
	user := User{Name: "Alice", Email: "alice@example.com"}
	s.db.Create(&user)

	s.NoError(s.db.Model(&user).Update("Email", "alice@newmail.com").Error)

	var updated User
	s.db.First(&updated, user.ID)
	s.Equal("alice@newmail.com", updated.Email)
}

func (s *mysqlTestContainerTestSuite) Test_Delete() {
	user := User{Name: "Bob", Email: "bob@example.com"}
	s.db.Create(&user)

	s.NoError(s.db.Delete(&user).Error)

	var found User
	err := s.db.First(&found, user.ID).Error
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}
