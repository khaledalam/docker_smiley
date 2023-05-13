package entities

import (
	"os"
)

type Process struct {
	os.Process
	Meta string `json:"Meta"`
}
