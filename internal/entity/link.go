package entity

import (
	"net/url"
)

type link struct {
	id     int
	alias  string
	url    url.URL
	userId int
}
