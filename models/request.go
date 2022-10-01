package models

type Request struct {
	Command string      `json:"command"`
	Param   interface{} `json:"param"`
	Data    interface{} `json:"data"`
}

type StreamSubject struct {
	StreamName     string
	StreamSubjects string
	SubjectName    string
}
