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
	
}

func (opt *Opt) ValidateStrategy() (err error){
	isValidStrategy :=  false;
	validStrategies := []string{"run","test"};
	for _,strategy :=  range validStrategies{
		if strategy == opt.Strategy {
			isValidStrategy = true;
		}
	}
}

type SandBoxInterface interface{
	Run() io.Reader
}

// SandBox executes code in sa
type SandBox struct{
	Opt Opt
}

func (s SandBox) Run() io.Reader{
	r := fmt.Sprintf("Language: %s\nStrategy: %s\nSource: %s\n",s.Opt.Language,s.Opt.Strategy,s.Opt.Source);
	return strings.NewReader(r);
}