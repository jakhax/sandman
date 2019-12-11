package runners;
import (
	"io"
)

const (
	HomePath = "/home/appuser/"
)


// Output : runner return vaue
type Output struct {
	StdOut  io.Reader
	StdErr io.Reader
}

// Runner interface
type Runner interface{
	Run(opt *Opt) (output Output, err error);
}

// CreateRunner returns a language specific runner
func CreateRunner(language string)(runner Runner, err error){
	return;
}


func SetupFromOpt(opt *Opt) (err error){
	opt.Dir = HomePath
	// get strategy
	if opt.Fixture != nil {
		opt.Strategy = "test"
	}else{
		opt.Strategy = "run"
	}
	// languages conf
	lc, err :=  GetLanguagesConf();
	if err != nil {
		return
	}
	opt.LanguagesConf = lc 
	


}


