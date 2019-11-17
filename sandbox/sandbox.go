package sandbox;
import (
	"io"
	"strings"
	"fmt"
)


type Opt struct{
	Strategy string `json:"strategy" validate:`
}

func (opt *Opt) OK() (err error){
	validStrategies := []string{"run","test"};
	isValidStrategy
}


type SandBoxInterface interface{
	Run() io.Reader
}

type SandBox struct{
	Opt Opt
}

func (s SandBox) Run() io.Reader{
	r := fmt.Sprintf("Language: %s\nStrategy: %s\nSource: %s\n",s.Opt.Language,s.Opt.Strategy,s.Opt.Source);
	return strings.NewReader(r);
}