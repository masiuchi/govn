package govn

import (
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Store struct {
	*Settings
}

func NewStore(settings *Settings) *Store {
	s := new(Store)
	s.Settings = settings
	return s
}

func (s *Store) GetValues(url string) *simplejson.Json {
	url = regexp.MustCompile("/").ReplaceAllString(url, "")
	url = s.Settings.ApiUrl + "?token=" + s.Settings.UserToken + "&url=" + url

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	json, _ := simplejson.NewJson(byteArray)
	return json
}
