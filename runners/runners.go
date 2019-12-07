package runners;
import (
	"io"
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