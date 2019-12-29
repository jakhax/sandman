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


var lConf []byte = []byte(
`
{
    "timeouts":{
        "default": 12000,
        "setup-shell":10000,
        "go": 15000,
        "haskell": 15000,
        "sql": 14000,
        "java": 20000,
        "kotlin": 23000,
        "groovy": 23000,
        "scala": 27000,
        "solidity": 20000
    },
    "fileExtensions":{
        "c": "c",
        "cpp": "cpp",
        "csharp": "cs",
        "elixir": "ex",
        "go": "go",
        "haskell": "hs",
        "java": "java",
        "javascript": "js",
        "julia": "jl",
        "kotlin": "kt",
        "ruby": "rb",
        "scala": "scala",
        "shell": "sh",
        "sql": "sql",
        "swift": "swift"
    },
    "images":{
        "base":"sandman/base",
        "python":"sandman/python-runner"
    },
    "workingDir":{
        "default":"/home/appuser"
    }
}`)

// GetLanguagesConf LanguagesConf
func GetLanguagesConf()(conf *LanguagesConf, err error){
	conf = &LanguagesConf{};
	err = json.Unmarshal(lConf,conf)
	return
}

//GetLanguagesConfFromFile LanguagesConf
func GetLanguagesConfFromFile()(conf *LanguagesConf, err error){
	data,err :=  ioutil.ReadFile("language_conf.json")
	if err != nil {
		return
	}
	conf = &LanguagesConf{};
	err = json.Unmarshal(data,conf)
	return
}