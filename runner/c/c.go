package c;

import (
	"fmt"
	"github.com/jakhax/sandman/runneropt"
	"io"
	"os"

	"github.com/jakhax/sandman/spawn"
	"io/ioutil"
)

// Runner code runner
type Runner struct {
}

//Setup method
func (r *Runner) Setup(opt *runneropt.Opt) (err error) {
	return
}

//Files method
func (r *Runner) Files(opt *runneropt.Opt) (err error) {
	return
}

//TestIntegration method
func (r *Runner) TestIntegration(opt *runneropt.Opt) (stdout, stderr io.Reader, err error) {
	return 
}


//SolutionOnly method
func (r *Runner) SolutionOnly(opt *runneropt.Opt) (stdout, stderr io.Reader, err error) {

	return
}

//SanitizeStdErr method
func (r *Runner) SanitizeStdErr(stderr io.Reader) (sanStderr io.Reader, err error) {
	return stderr, err
}

//SanitizeStdOut method
func (r *Runner) SanitizeStdOut(stdout io.Reader) (sanStdout io.Reader, err error) {
	return stdout, err
}

//TransformOutput method
func (r *Runner) TransformOutput(stdout, stderr io.Reader) (tStdout, tStderr io.Reader, err error) {
	return stdout, stderr, err
}
