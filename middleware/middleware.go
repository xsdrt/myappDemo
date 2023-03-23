package middleware

import (
	"myappDemo/data"
	"myappDemo/vendor/github.com/xsdrt/hiSpeed"
)

type Middleware struct {
	App    *hiSpeed.HiSpeed
	Models data.Models
}
