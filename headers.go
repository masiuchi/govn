package govn

import (
	"net/http"
	"regexp"
	"strings"
)

type Headers struct {
	UnmaskedUrl      string
	Url              string
	Protocol         string
	UnmaskedHost     string
	Host             string
	UnmaskedPathName string
	PathName         string
	RedisUrl         string
	PathLang         string
	*Settings
	Request     *http.Request
	BrowserLang string
}

func NewHeaders(req *http.Request, settings *Settings) *Headers {
	h := new(Headers)
	h.Request = req
	h.Settings = settings
	h.Protocol = h.Request.URL.Scheme
	h.UnmaskedHost = h.Request.Host
	uri := h.Request.RequestURI
	if uri != "" {
		if regexp.MustCompile(`^[^/]`).MatchString(h.Request.URL.Path) {
			h.Request.RequestURI = "/"
		}
		//		h.Request.RequestURI += h.Request.URL.Path
		if len(h.Request.URL.RawQuery) > 0 {
			h.Request.RequestURI += "?" + h.Request.URL.RawQuery
		}
	}
	if regexp.MustCompile(`://`).MatchString(uri) {
		h.Request.RequestURI = regexp.MustCompile(`^.*://[^/]+`).ReplaceAllString(uri, "")
	}
	split := strings.Split(h.Request.RequestURI, "?")
	h.UnmaskedPathName = split[0]
	if !regexp.MustCompile("/$").MatchString(h.UnmaskedPathName) {
		if !regexp.MustCompile(`/[^/.]+\.[^/.]+$`).MatchString(h.UnmaskedPathName) {
			h.UnmaskedPathName += "/"
		}
	}
	h.UnmaskedUrl = h.Protocol + "://" + h.UnmaskedHost + h.UnmaskedPathName
	if settings.UrlPattern == "subdomain" {
		h.Host = h.RemoveLang(h.Request.Host, h.LangCode())
	} else {
		h.Host = h.Request.Host
	}

	return h
}

func (h *Headers) LangCode() string {
	if h.PathLang != "" && len(h.PathLang) > 0 {
		return h.GetPathLang()
	} else {
		return h.Settings.DefaultLang
	}
}

func (h *Headers) GetPathLang() string {
	if h.PathLang == "" {
		r := h.Settings.UrlPatternReg
		match := r.FindStringSubmatch(h.Request.URL.Host + h.Request.RequestURI)
		result := make(map[string]string)
		if match != nil {
			for i, name := range r.SubexpNames() {
				// result[name] = match[i]
				if i != 0 {
					result[name] = match[i]
				}
			}
			if lang, ok := result["lang"]; ok {
				if GetLang(lang) != nil {
					h.PathLang = GetCode(lang)
				}
			}
		}
	}
	return h.PathLang
}

func (h *Headers) GetBrowserLang() string {
	if h.BrowserLang == "" {
		lang, err := h.Request.Cookie("wovn_selected_lang")
		if err != nil && GetLang(lang.Value) != nil {
			h.BrowserLang = lang.Value
		} else {
			r := regexp.MustCompile(`[,;]`)
			acceptLangs := r.Split(h.Request.Header.Get("Accept-Language"), -1)
			for _, l := range acceptLangs {
				if GetLang(l) != nil {
					h.BrowserLang = l
					break
				}
			}
		}
	}
	return h.BrowserLang
}

func (h *Headers) Redirect(lang string) *http.Header {
	if lang == "" {
		lang = h.GetBrowserLang()
	}
	redirectHeader := new(http.Header)
	redirectHeader.Set("location", h.RedirectLocation(lang))
	redirectHeader.Set("content-length", "0")
	return redirectHeader
}

func (h *Headers) RedirectLocation(lang string) string {
	if lang == h.Settings.DefaultLang {
		return h.Protocol + "://" + h.Url
	} else {
		location := h.Url
		switch h.Settings.UrlPattern {
		case "query":
			r := regexp.MustCompile(`?`)
			if !r.MatchString(location) {
				location += "?wovn=" + lang
			} else {
				r = regexp.MustCompile(`(?|&)wovn=`)
				if !r.MatchString(h.Request.RequestURI) {
					location += "&wovn=" + lang
				}
			}
		case "subdomain":
			location = strings.ToLower(lang) + "." + location
		default:
			r := regexp.MustCompile("(/|$)")
			location = r.ReplaceAllString(location, "/"+lang+"/")
		}
		return h.Protocol + "://" + location
	}
}

func (h *Headers) RequestOut(defLang string) *http.Request {
	if defLang == "" {
		defLang = h.Settings.DefaultLang
	}
	switch h.Settings.UrlPattern {
	case "query":
		if h.Request.RequestURI != "" {
			h.Request.RequestURI = h.RemoveLang(h.Request.RequestURI)
		}
		if h.Request.URL.RawQuery != "" {
			h.Request.URL.RawQuery = h.RemoveLang(h.Request.URL.RawQuery)
		}
		if h.Request.URL.RawPath != "" {
			h.Request.URL.RawPath = h.RemoveLang(h.Request.URL.RawPath)
		}
	case "subdomain":
		h.Request.Host = h.RemoveLang(h.Request.Host)
		/*
			if _, ok := h.Env["HTTP_REFERER"]; ok {
				h.Env["HTTP_REFERER"] = h.RemoveLang(h.Env["HTTP_REFERER"])
			}
		*/
	default:
		h.Request.RequestURI = h.RemoveLang(h.Request.RequestURI)
		if h.Request.URL.RawPath != "" {
			h.Request.URL.RawPath = h.RemoveLang(h.Request.URL.RawPath)
		}
		//		h.Env["PATH_INFO"] = h.RemoveLang(h.Env["PATH_INFO"])
		if h.Request.URL.RawPath != "" {
			h.Request.URL.RawPath = h.RemoveLang(h.Request.URL.RawPath)
		}
	}

	return h.Request
}

func (h *Headers) RemoveLang(args ...string) string {
	if len(args) != 1 && len(args) != 2 {
		panic("Invalid arguments") // TODO
	}

	uri := args[0]
	var lang string
	if len(args) == 2 {
		lang = args[1]
	} else {
		lang = h.GetPathLang()
	}

	switch h.Settings.UrlPattern {
	case "query":
		r := regexp.MustCompile(`(^|\?|&)wovn=` + lang + `(&|$)`)
		uri = r.ReplaceAllString(uri, "$1")
		r = regexp.MustCompile(`(\?|&)$`)
		uri = r.ReplaceAllString(uri, "")
	case "subdomain":
		uri = strings.ToLower(uri)
		r := regexp.MustCompile(`(^|(//))` + lang + `\.`)
		uri = r.ReplaceAllString(uri, "$1")
	default:
		r := regexp.MustCompile(`//` + lang + `(/|$)`)
		uri = r.ReplaceAllString(uri, "/")
	}
	return uri
}

func (h *Headers) Out(req *http.Request) *http.Request {
	r := regexp.MustCompile(`//` + h.Host)
	l := req.Header.Get("Location")
	if r.MatchString(l) {
		switch h.Settings.UrlPattern {
		case "query":
			r = regexp.MustCompile(`\?`)
			if r.MatchString(l) {
				l += "&"
			} else {
				l += "?"
			}
			req.Header.Set("Location", l+h.LangCode())
		case "subdomain":
			r = regexp.MustCompile(`//([^.]+)`)
			l = r.ReplaceAllString(l, "//"+h.LangCode()+".$1")
			req.Header.Set("Location", l)
		default:
			r = regexp.MustCompile(`(//[^/]+)`)
			l = r.ReplaceAllString(l, "$1/"+h.LangCode())
			req.Header.Set("Location", l)
		}
	}
	return req
}
