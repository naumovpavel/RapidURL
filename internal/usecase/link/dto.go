package link

import (
	"net/url"
)

type SaveLinkDTO struct {
	Alias  string  `json:"alias,omitempty"`
	Url    url.URL `json:"url"`
	UserId int     `json:"userId"`
}
