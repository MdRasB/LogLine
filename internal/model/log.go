package model

type Logs struct {
	Level string
	Message string
	Service string
	Timestamp string
	Metadata map[string]interface{}
}