package middleware

import (
	"myappDemo/data"

	"github.com/xsdrt/hiSpeed"
)

type Middleware struct {
	App    *hiSpeed.HiSpeed
	Models data.Models
}
