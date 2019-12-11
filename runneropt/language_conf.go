package runneropt

import (
	"io/ioutil"
	"encoding/json"
)

// LanguagesConf conf
type LanguagesConf struct{
	FileExtensions map[string]string
	Timeouts map[string]int
	Images map[string]string
	WorkingDir map[string]string
}

// GetLanguagesConf LanguagesConf
func GetLanguagesConf()(conf *LanguagesConf, err error){
	data,err :=  ioutil.ReadFile("language_conf.json")
	if err != nil {
		return
	}
	conf = &LanguagesConf{};
	err = json.Unmarshal(data,conf)
	return
}