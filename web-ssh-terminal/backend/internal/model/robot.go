package model

type AuthType string

const (
	AuthPassword   AuthType = "password"
	AuthPrivateKey AuthType = "privateKey"
)

type Robot struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Host        string   `json:"host" yaml:"host"`
	Port        int      `json:"port" yaml:"port"`
	Username    string   `json:"username" yaml:"username"`
	AuthType    AuthType `json:"authType" yaml:"authType"`
	Description string   `json:"description,omitempty" yaml:"description"`
	IsOnline    bool     `json:"isOnline" yaml:"-"`
}
