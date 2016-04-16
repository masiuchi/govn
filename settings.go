package govn

import (
	"regexp"
)

type Settings struct {
	UserToken      string
	SecretKey      string
	UrlPattern     string
	UrlPatternReg  *regexp.Regexp
	Query          []string
	ApiUrl         string
	DefaultLang    string
	SupportedLangs []string
	TestMode       bool
	TestUrl        string
}

const (
	DefaultUrlPattern  = "path"
	DefaultApiUrl      = "https://api.wovn.io/v0/values"
	DefaultDefaultLang = "en"
)

var DefaultUrlPatternReg = regexp.MustCompile("/(?P<lang>[^/.?]+)")

func NewSettings() *Settings {
	s := new(Settings)

	s.UrlPattern = DefaultUrlPattern
	s.UrlPatternReg = DefaultUrlPatternReg
	s.ApiUrl = DefaultApiUrl
	s.DefaultLang = DefaultDefaultLang
	s.SupportedLangs = append(s.SupportedLangs, DefaultDefaultLang)

	return s
}

func (s *Settings) IsValid() bool {
	valid := true
	errors := []string{}

	if len(s.UserToken) < 5 || len(s.UserToken) > 6 {
		valid = false
		errors = append(errors, "User token "+s.UserToken+" is not valid.")
	}
	if s.SecretKey == "" {
		valid = false
		errors = append(errors, "Secret key must be required.")
	}
	if s.UrlPattern == "" {
		valid = false
		errors = append(errors, "Url pattern must be required.")
	}
	if s.ApiUrl == "" {
		valid = false
		errors = append(errors, "API url must be required.")
	}
	if s.DefaultLang == "" {
		valid = false
		errors = append(errors, "Default lang must be required.")
	}
	if len(s.SupportedLangs) < 1 {
		valid = false
		errors = append(errors, "Supported langs must be required.")
	}

	return valid
}

func (s *Settings) Initialize() {
	s.DefaultLang = GetCode(s.DefaultLang)
	if len(s.SupportedLangs) == 0 {
		s.SupportedLangs = append(s.SupportedLangs, s.DefaultLang)
	}
	switch s.UrlPattern {
	case "path":
		s.UrlPatternReg = regexp.MustCompile("(?P<lang>[^/.?]+)")
	case "query":
		s.UrlPatternReg = regexp.MustCompile("((\\?.*&)|\\?)wovn=(?P<lang>[^&]+)(&|$)")
	default:
		s.UrlPatternReg = regexp.MustCompile(`^(?P<lang>[^.]+)\.`)
	}
}
