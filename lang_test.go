package govn

import (
	"testing"
)

func TestLangsExists(t *testing.T) {
	if LANG == nil {
		t.Errorf("Lang should not be nil")
	}
}

func TestLangsSize(t *testing.T) {
	if len(LANG) != 28 {
		t.Errorf("len(LANG) should be 28: %v", len(LANG))
	}
}

func TestKeysExists(t *testing.T) {
	for key, value := range LANG {
		if value.Name == "" {
			t.Errorf("Name properties is empty.\n(key: %v, value: %v)", key, value)
		}
		if value.Code == "" {
			t.Errorf("Code properties is empty.\n(key: %v, value: %v)", key, value)
		}
		if value.EnglishName == "" {
			t.Errorf("EnglishName properties is empty.\n(key: %v, value: %v)", key, value)
		}

		actual := value.Code
		expected := key
		if actual != expected {
			t.Errorf("got %v\nwant %v", actual, expected)
		}
	}
}

func TestGetCodeWithValidCode(t *testing.T) {
	actual := GetCode("ms")
	expected := "ms"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGetCodeWithCapitalLetters(t *testing.T) {
	actual := GetCode("zh-cht")
	expected := "zh-CHT"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGetCodeWithValidEnglishName(t *testing.T) {
	actual := GetCode("Portuguese")
	expected := "pt"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGetCodeWithValidNativeName(t *testing.T) {
	actual := GetCode("हिन्  दी")
	expected := "hi"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGetCodeWithInvalidName(t *testing.T) {
	actual := GetCode("WOVN4LYFE")
	expected := ""
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestGetCodeWithEmptyString(t *testing.T) {
	actual := GetCode("")
	expected := ""
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
