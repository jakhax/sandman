package runners;
import (
	"io"
	"io/ioutil"
	"encoding/json"
	"github.com/jakhax/sandman/utils"
)

// Opt option for running code in container
type Opt struct{
	Language string `json:"language"`
	Code string `json:"code"`
	Fixture string `json:"fixture"`
	SetupCode string `json:"setup_code"`
	TestFramework string `json:"test_framework"`
	Format string `json:"format"`
	Shell string `json:"shell"`
	Timeout int `json:"timeout"`
	Memory int `json:"memory"`
	CPU int `json:"cpu"`
}

// OK validates the Opt struct
func (opt *Opt) OK() (err error){

	if opt.Code == ""{
		err = utils.ValidationError{Message:"Code is required"}
		return
	}
	if opt.Language == ""{
		err = utils.ValidationError{Message:"Language is required"}
		return
	}
	return;
}

// NewOptFromReader creates Opt
func NewOptFromReader(optR io.Reader)(opt *Opt, err error){
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