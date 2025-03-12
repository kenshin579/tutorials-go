package model

type Student struct {
	Name string
}

type StudentResponse struct {
	Data        interface{}       `json:"data,omitempty"`
	Message     string            `json:"message,omitempty"`
	MessageList []string          `json:"messageList,omitempty"`
	Errors      map[string]string `json:"errors,omitempty"`
	Student     Student           `json:"student,omitempty"`
	StudentList []Student         `json:"studentList,omitempty"`
}
