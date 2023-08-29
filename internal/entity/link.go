package entity

import (
	"net/url"
)

type Link struct {
	Id     int
	Alias  string
	Url    *url.URL
	UserId int
}
