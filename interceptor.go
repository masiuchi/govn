package govn

import (
	"github.com/bitly/go-simplejson"
	//	"gopkg.in/xmlpath.v2"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Interceptor struct {
	*Store
	HttpRequest *http.Request
}

func NewInterceptor(settings *Settings) *Interceptor {
	interceptor := new(Interceptor)
	settings.Initialize()
	interceptor.Store = NewStore(settings)
	return interceptor
}

func (interceptor *Interceptor) Call(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !interceptor.Store.Settings.IsValid() {
			next.ServeHTTP(w, r)
			return
		}

		interceptor.HttpRequest = r
		headers := NewHeaders(r, interceptor.Store.Settings)

		if interceptor.Store.Settings.TestMode && interceptor.Store.Settings.TestUrl != headers.Url {
			next.ServeHTTP(w, r)
			return
		}

		if headers.GetPathLang() == interceptor.Store.Settings.DefaultLang {
			redirectUrl := headers.RedirectLocation(interceptor.Store.Settings.DefaultLang)
			http.Redirect(w, r, redirectUrl, 307)
			return
		}

		ww := NewWrapperWriter(w)
		lang := headers.LangCode()
		next.ServeHTTP(ww, headers.RequestOut(lang))

		if regexp.MustCompile("html").MatchString(ww.Header().Get("Content-Type")) {
			values := interceptor.Store.GetValues(headers.RedisUrl)

			u := map[string]string{
				"protocol": headers.Protocol,
				"host":     headers.Host,
				"pathname": headers.PathName,
			}

			if !(ww.Status >= 100 && ww.Status < 200) && ww.Status != 302 {
				ww.Body = interceptor.SwitchLang(ww.Body, values, u, lang, headers)
			}
		}

		ww.Flush()
	})
}

func AddLangCode(href, pattern, lang string, h *Headers) string {
	if regexp.MustCompile("^(#.*)?$").MatchString(href) {
		return href
	}
	newHref := href
	if len(href) > 0 && regexp.MustCompile("^(https?:)?//").MatchString(strings.ToLower(href)) {
		uri, err := url.Parse(href)
		if err != nil {
			return newHref
		}
		if strings.ToLower(uri.Host) == strings.ToLower(h.Host) {
			switch pattern {
			case "subdomain":
				subDomain := regexp.MustCompile(`\/\/([^\.]*)\.`).FindString(href)
				subCode := GetCode(subDomain)
				if len(subCode) > 0 && strings.ToLower(subCode) == strings.ToLower(lang) {
					newHref = regexp.MustCompile("(?i)"+lang).ReplaceAllString(href, strings.ToLower(lang))
				} else {
					newHref = regexp.MustCompile(`(\/\/)([^\.]*)`).ReplaceAllString(href, "$1"+strings.ToLower(lang)+".$2")
				}
			case "query":
				if regexp.MustCompile(`?`).MatchString(href) {
					newHref = href + "&wovn=" + lang
				} else {
					newHref = href + "?wovn=" + lang
				}
			default:
				newHref = regexp.MustCompile(`([^\\.]*\\.[^/]*)(/|$)`).ReplaceAllString(href, "$1"+lang+"/")
			}
		}
	} else if len(href) > 0 {
		switch pattern {
		case "subdomain":
			langUrl := h.Protocol + "://" + strings.ToLower(lang) + "." + h.Host
			currentDir := regexp.MustCompile(`[^\/]*\.[^\.]{2,6}$`).ReplaceAllString(h.PathName, "")
			if regexp.MustCompile(`^\.\..*$`).MatchString(href) {
				newHref = langUrl + "/" + regexp.MustCompile(`^\.\.\/`).ReplaceAllString(href, "")
			} else if regexp.MustCompile(`^\..*$`).MatchString(href) {
				newHref = langUrl + currentDir + "/"
				newHref += regexp.MustCompile(`^\.\/`).ReplaceAllString(href, "")
			} else if regexp.MustCompile(`^\/.*$`).MatchString(href) {
				newHref = langUrl + href
			} else {
				newHref = langUrl + currentDir + "/" + href
			}
		case "query":
			if regexp.MustCompile(`?`).MatchString(href) {
				newHref = href + "&wovn=" + lang
			} else {
				newHref = href + "?wovn=" + lang
			}
		default:
			if regexp.MustCompile("^/").MatchString(href) {
				newHref = "/" + lang + href
			} else {
				currentDir := regexp.MustCompile(`[^\/]*\.[^\.]{2,6}$`).ReplaceAllString(h.PathName, "")
				newHref = "/" + lang + currentDir + href
			}
		}
	}
	return newHref
}

func CheckWovnIgnore(n xml.Node) bool {
	if n.Attr("wovn-ignore") == "" {
		return true
	} else if n.Name() == "html" {
		return false
	}
	return CheckWovnIgnore(n.Parent())
}

func (interceptor *Interceptor) SwitchLang(body string, values *simplejson.Json, u map[string]string, lang string, headers *Headers) string {

	if lang == "" {
		lang = interceptor.Store.Settings.DefaultLang
	}
	lang = GetCode(lang)
	textIndex := values.Get("text_vals")
	srcIndex := values.Get("img_vals")
	imgSrcPrefix, _ := values.Get("img_src_prefix").String()

	doc, err := gokogiri.ParseHtml([]byte(body))
	if err != nil {
		return body // Do nothing.
	}
	defer doc.Free()

	x := xpath.Compile("//html[@wovn-ignore]")
	nodes, err := doc.Search(x)
	if err != nil {
	} else if len(nodes) > 0 {
		h, _ := doc.ToHtml(nil, nil)
		html := string(h)
		re := regexp.MustCompile(`href="([^"]*)"`)
		html = re.ReplaceAllStringFunc(html, func(s string) string {
			matched := re.FindStringSubmatch(s)
			replaced, err := url.QueryUnescape(matched[0])
			if err != nil {
				return s
			} else {
				return replaced
			}
		})
		return html
	}

	if lang != interceptor.Store.Settings.DefaultLang {
		pattern := interceptor.Store.Settings.UrlPattern

		x = xpath.Compile("//a")
		nodes, err = doc.Search(x)
		if err != nil {
		} else {
			for _, n := range nodes {
				if !CheckWovnIgnore(n) {
					continue
				}
				href := n.Attr("href")
				newHref := AddLangCode(href, pattern, lang, headers)
				n.SetAttr("href", newHref)
			}
		}

		x = xpath.Compile("//form")
		nodes, err = doc.Search(x)
		if err != nil {
		} else {
			for _, n := range nodes {
				if !CheckWovnIgnore(n) {
					continue
				}
				method := n.Attr("method")
				if pattern == "query" && (method == "" || strings.ToUpper(method) == "GET") {
					input := doc.CreateElementNode("input")
					input.SetAttr("type", "hidden")
					input.SetAttr("name", "wovn")
					input.SetAttr("value", lang)
					if n.CountChildren() > 0 {
						n.FirstChild().AddPreviousSibling(input)
					} else {
						n.AddChild(input)
					}
				} else {
					action := n.Attr("action")
					newAction := AddLangCode(action, pattern, lang, headers)
					n.SetAttr("action", newAction)
				}
			}
		}
	}

	x = xpath.Compile("//text()")
	nodes, err = doc.Search(x)
	if err != nil {
	} else {
		for _, n := range nodes {
			if !CheckWovnIgnore(n) {
				continue
			}
			nodeText := regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(n.Content(), "")
			if srcs, ok := textIndex.CheckGet(nodeText); ok {
				if langs, ok := srcs.CheckGet(lang); ok {
					arr, err := langs.Array()
					if err == nil && len(arr) > 0 {
						l := langs.GetIndex(0)
						data, err := l.Get("data").String()
						if err == nil {
							content := strings.Replace(n.Content(), nodeText, data, -1)
							n.SetContent(content)
						}
					}
				}
			}
		}
	}

	x = xpath.Compile("//meta")
	nodes, err = doc.Search(x)
	if err == nil {
		for _, n := range nodes {
			if !CheckWovnIgnore(n) {
				continue
			}
			attr := n.Attr("name")
			if attr == "" {
				attr = n.Attr("property")
			}
			if !regexp.MustCompile(`^(description|title|og:title|og:description|twitter:title|twitter:description)$`).MatchString(attr) {
				continue
			}
			content := n.Attr("content")
			content = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(content, "")
			if srcs, ok := textIndex.CheckGet(content); ok {
				if langs, ok := srcs.CheckGet(lang); ok {
					arr, err := langs.Array()
					if err == nil && len(arr) > 0 {
						l := langs.GetIndex(0)
						data, err := l.Get("data").String()
						if err == nil {
							v := strings.Replace(n.Attr("content"), content, data, -1)
							n.SetAttr("content", v)
						}
					}
				}
			}
		}
	}

	x = xpath.Compile("//img")
	nodes, err = doc.Search(x)
	if err == nil {
		for _, n := range nodes {
			if !CheckWovnIgnore(n) {
				continue
			}
			h, _ := n.ToHtml(nil, nil)
			html := string(h)
			matched := regexp.MustCompile(`(?i)src=['"]([^'"]*)['"]`).FindStringSubmatch(html)
			if len(matched) > 0 {
				src := matched[1]
				if !regexp.MustCompile(`://`).MatchString(src) {
					if regexp.MustCompile(`^/`).MatchString(src) {
						src = u["protocol"] + "://" + u["host"] + src
					} else {
						src = u["protocol"] + "://" + u["host"] + u["path"] + src
					}
				}

				if srcs, ok := srcIndex.CheckGet(src); ok {
					if langs, ok := srcs.CheckGet(lang); ok {
						arr, err := langs.Array()
						if err == nil && len(arr) > 0 {
							l := langs.GetIndex(0)
							data, err := l.Get("data").String()
							if err == nil {
								n.SetAttr("src", imgSrcPrefix+data)
							}
						}
					}
				}
			}

			alt := n.Attr("alt")
			if alt != "" {
				alt = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(alt, "")
				if srcs, ok := textIndex.CheckGet(alt); ok {
					if langs, ok := srcs.CheckGet(lang); ok {
						arr, err := langs.Array()
						if err == nil && len(arr) > 0 {
							l := langs.GetIndex(0)
							data, err := l.Get("data").String()
							if err == nil {
								v := strings.Replace(n.Attr("alt"), alt, data, -1)
								n.SetAttr("alt", v)
							}
						}
					}
				}
			}
		}
	}

	x = xpath.Compile("//script")
	nodes, err = doc.Search(x)
	if err == nil {
		for _, n := range nodes {
			src := n.Attr("src")
			if regexp.MustCompile(`//j.(dev-)?wovn.io(:3000)?/`).MatchString(src) {
				n.Remove()
			}
		}
	}

	var parentNode xml.Node
	x = xpath.Compile("head")
	nodes, err = doc.Search(x)
	if err == nil && len(nodes) > 0 {
		parentNode = nodes[0]
	}
	if parentNode == nil {
		x = xpath.Compile("body")
		nodes, err := doc.Search(x)
		if err == nil && len(nodes) > 0 {
			parentNode = nodes[0]
		}
	}
	if parentNode == nil {
		parentNode = doc
	}

	insertNode := doc.CreateElementNode("script")
	insertNode.SetAttr("src", "//j.wovn.io/1")
	insertNode.SetAttr("async", "")
	version := ""
	dataWovnio := "key=" + interceptor.Store.Settings.UserToken
	dataWovnio += "&backend=true&currentLang=" + lang
	dataWovnio += "&defaultLang=" + interceptor.Store.Settings.DefaultLang
	dataWovnio += "&urlPattern=" + interceptor.Store.Settings.UrlPattern
	dataWovnio += "&version=" + version
	insertNode.SetAttr("data-wovnio", dataWovnio)
	insertNode.SetContent(" ")

	if parentNode.CountChildren() > 0 {
		parentNode.FirstChild().AddPreviousSibling(insertNode)
	} else {
		parentNode.AddChild(insertNode)
	}

	publishedLangs := GetLangs(values)
	for _, l := range publishedLangs {
		insertNode = doc.CreateElementNode("link")
		insertNode.SetAttr("rel", "alternate")
		insertNode.SetAttr("hreflang", l)
		insertNode.SetAttr("href", headers.RedirectLocation(l))
		parentNode.AddChild(insertNode)
	}

	x = xpath.Compile("html")
	nodes, err = doc.Search(x)
	if err != nil && len(nodes) == 0 {
		x = xpath.Compile("HTML")
		nodes, err = doc.Search(x)
	}
	if nodes != nil && len(nodes) > 0 {
		nodes[0].SetAttr("lang", lang)
	}

	h, _ := doc.ToHtml(nil, nil)
	html := string(h)
	re := regexp.MustCompile(`href="([^"]*)"`)
	html = re.ReplaceAllStringFunc(html, func(s string) string {
		matched := re.FindStringSubmatch(s)
		replaced, err := url.QueryUnescape(matched[0])
		if err != nil {
			return s
		} else {
			return replaced
		}
	})

	return html
}

func GetLangs(values *simplejson.Json) []string {
	langs := []string{}

	textVals, ok := values.CheckGet("text_vals")
	if ok {
		m, err := textVals.Map()
		if err != nil {
			for key, _ := range m {
				langs = append(langs, key)
			}
		}
	}

	imgVals, ok := values.CheckGet("img_vals")
	if ok {
		m, err := imgVals.Map()
		if err != nil {
			for key, _ := range m {
				langs = append(langs, key)
			}
		}
	}

	return langs
}
