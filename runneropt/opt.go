package runneropt;
import (
	"io"
	"io/ioutil"
	"encoding/json"
	"github.com/jakhax/sandman/utils"
)

// Opt option for running code in container
type Opt struct{
	Strategy string 
	Language string `json:"language"`
	Code []byte `json:"code"`
	Fixture []byte `json:"fixture"`
	SetupCode []byte `json:"setup_code"`
	TestFramework string `json:"test_framework"`
	Format string `json:"format"`
	Shell []byte `json:"shell"`
	EntyFile string `json:"entry_file"`
	ProjectMode bool `json:"project_mode"`
	Files []File `json:"files"`
	Timeout int `json:"timeout"`
	Memory int `json:"memory"`
	CPU int `json:"cpu"`
	// working dir
	Dir string `json:"dir"`
	GitURL string `json:"git_url"`
	Env []string
	LanguagesConf *LanguagesConf
}

// File for project mode
type File struct{
	Name string
	Content []byte
}

// OK validates the Opt struct
func (opt *Opt) OK() (err error){

	if opt.Code == nil{
		err = utils.ValidationError{Message:"Code is required"}
		return
	}
	if opt.Language == ""{
		err = utils.ValidationError{Message:"Language is required"}
		return
	}
	return;
}

// NewOptFromJSON creates Opt
func NewOptFromJSON(optR io.Reader)(opt *Opt, err error){
	optData,err := ioutil.ReadAll(optR);
	if err != nil{
		return;
	}
	opt = &Opt{}
	err = json.Unmarshal(optData,opt);
	if err != nil{
		return;
	} 
	err = opt.OK()
	return;
}