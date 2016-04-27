package govn

import (
	"strings"
)

type Lang struct {
	Name        string
	Code        string
	EnglishName string
}

var LANG = map[string]Lang{
	"ar": Lang{
		Name:        "ﺎﻠﻋﺮﺒﻳﺓ",
		Code:        "ar",
		EnglishName: "Arabic",
	},
	"bg": Lang{
		Name:        "Български",
		Code:        "bg",
		EnglishName: "Bulgarian",
	},
	"zh-CHS": Lang{
		Name:        "简体中文",
		Code:        "zh-CHS",
		EnglishName: "Simp Chinese",
	},
	"zh-CHT": Lang{
		Name:        "繁體中文",
		Code:        "zh-CHT",
		EnglishName: "Trad Chinese",
	},
	"da": Lang{
		Name:        "Dansk",
		Code:        "da",
		EnglishName: "Danish",
	},
	"nl": Lang{
		Name:        "Nederlands",
		Code:        "nl",
		EnglishName: "Dutch",
	},
	"en": Lang{
		Name:        "English",
		Code:        "en",
		EnglishName: "English",
	},
	"fi": Lang{
		Name:        "Suomi",
		Code:        "fi",
		EnglishName: "Finnish",
	},
	"fr": Lang{
		Name:        "Français",
		Code:        "fr",
		EnglishName: "French",
	},
	"de": Lang{
		Name:        "Deutsch",
		Code:        "de",
		EnglishName: "German",
	},
	"el": Lang{
		Name:        "Ελληνικά",
		Code:        "el",
		EnglishName: "Greek",
	},
	"he": Lang{
		Name:        "עברית",
		Code:        "he",
		EnglishName: "Hebrew",
	},
	"id": Lang{
		Name:        "Bahasa Indonesia",
		Code:        "id",
		EnglishName: "Indonesian",
	},
	"it": Lang{
		Name:        "Italiano",
		Code:        "it",
		EnglishName: "Italian",
	},
	"ja": Lang{
		Name:        "日本語",
		Code:        "ja",
		EnglishName: "Japanese",
	},
	"ko": Lang{
		Name:        "한국어",
		Code:        "ko",
		EnglishName: "Korean",
	},
	"ms": Lang{
		Name:        "Bahasa Melayu",
		Code:        "ms",
		EnglishName: "Malay",
	},
	"no": Lang{
		Name:        "Norsk",
		Code:        "no",
		EnglishName: "Norwegian",
	},
	"pl": Lang{
		Name:        "Polski",
		Code:        "pl",
		EnglishName: "Polish",
	},
	"pt": Lang{
		Name:        "Português",
		Code:        "pt",
		EnglishName: "Portuguese",
	},
	"ru": Lang{
		Name:        "Русский",
		Code:        "ru",
		EnglishName: "Russian",
	},
	"es": Lang{
		Name:        "Español",
		Code:        "es",
		EnglishName: "Spanish",
	},
	"sv": Lang{
		Name:        "Svensk",
		Code:        "sv",
		EnglishName: "Swedish",
	},
	"th": Lang{
		Name:        "ภาษาไทย",
		Code:        "th",
		EnglishName: "Thai",
	},
	"hi": Lang{
		Name:        "हिन्  दी",
		Code:        "hi",
		EnglishName: "Hindi",
	},
	"tr": Lang{
		Name:        "Türkçe",
		Code:        "tr",
		EnglishName: "Turkish",
	},
	"uk": Lang{
		Name:        "Українська",
		Code:        "uk",
		EnglishName: "Ukrainian",
	},
	"vi": Lang{
		Name:        "Tiếng Việt",
		Code:        "vi",
		EnglishName: "Vietnamese",
	},
}

func GetCode(langName string) string {
	if langName == "" {
		return ""
	}
	if _, ok := LANG[langName]; ok {
		return langName
	}
	for _, l := range LANG {
		lowerLangName := strings.ToLower(langName)
		if lowerLangName == strings.ToLower(l.Name) || lowerLangName == strings.ToLower(l.EnglishName) || lowerLangName == strings.ToLower(l.Code) {
			return l.Code
		}
	}
	return ""
}

func GetLang(langName string) *Lang {
	code := GetCode(langName)
	if code != "" {
		lang := LANG[code]
		return &lang
	} else {
		return nil
	}
}
