package python
import (
	"github.com/jakhax/sandman/runneropt"
	"io"
	"os"
	"fmt"

	"io/ioutil"
	"github.com/jakhax/sandman/spawn"
)

// Runner code runner
type Runner struct {

}

//Setup method
func (r *Runner) Setup(opt *runneropt.Opt) (err error){
	return
}

//Files method
func (r *Runner) Files(opt *runneropt.Opt)(err error){
	return
}

//TestIntegration method
func (r *Runner) TestIntegration(opt *runneropt.Opt)(stdout,stderr io.Reader, err error){
	return r.Python3Unittest(opt)
}

// Python3Unittest method
// .
// |-- setup.py
// |-- solution.py
// |-- codewars_testwrapper.py
// `-- test
//     |-- __init__.py
//     |-- __main__.py
//     `-- test_solution.py
// inspired by http://stackoverflow.com/a/27630375
func (r *Runner) Python3Unittest(opt *runneropt.Opt)(stdout,stderr io.Reader, err error){
	if len(opt.SetupCode) != 0{
		errX := ioutil.WriteFile("setup.py",opt.SetupCode,0644)
		if errX != nil{
			err = errX
			return
		}
		code := "from solution import *\n"+string(opt.Code)
		opt.Code = []byte(code)
	}
	err = ioutil.WriteFile("solution.py",opt.Code,0644)
	if err != nil{
		return
	}
	err = os.Mkdir("test",0777)
	if err != nil{
		return
	}
	err = ioutil.WriteFile("test/__init__.py",[]byte(""),0644)
	if err != nil{
		return
	}
	testLoader := fmt.Sprint(`import unittest
import timeout_decorator
from codewars_unittest import CodewarsTestRunner
def load_tests(loader, tests, pattern):
	return loader.discover('.')
unittest.main(testRunner=CodewarsTestRunner())
GLOBAL_TIMEOUT=3
timeout_decorator.timeout(GLOBAL_TIMEOUT)(unittest.main)(testRunner=CodewarsTestRunner())`)

	err = ioutil.WriteFile("test/__main__.py",[]byte(testLoader),0644);
	if err != nil{
		return
	}
	err = ioutil.WriteFile("test/test_solution.py",opt.Fixture,0644)
	if err != nil{
		return
	}
	spawnOpt := &spawn.SpawnOpt{
		Dir:opt.Dir,
		Timeout:opt.Timeout,
	};
	var stdin io.Reader
	stdout, stderr, err = spawn.Spwan(spawnOpt,"python3",[]string{"test"},stdin)
	return
}

//SolutionOnly method
func (r *Runner) SolutionOnly(opt *runneropt.Opt)(stdout,stderr io.Reader, err error){
	spawnOpt := &spawn.SpawnOpt{
		Dir:opt.Dir,
		Timeout:opt.Timeout,
	};
	code := string(opt.Code)
	var stdin io.Reader
	stdout, stderr, err = spawn.Spwan(spawnOpt,"python",[]string{"-c",code},stdin)
	return 
}

//SanitizeStdErr method
func (r *Runner) SanitizeStdErr(stderr io.Reader)(sanStderr io.Reader, err error){
	return stderr,err
}

//SanitizeStdOut method
func (r *Runner) SanitizeStdOut(stdout io.Reader)(sanStdout io.Reader, err error){
	return stdout,err
}

//TransformOutput method
func (r *Runner) TransformOutput(stdout,stderr io.Reader)(tStdout,tStderr io.Reader, err error){
	return stdout,stderr,err
}