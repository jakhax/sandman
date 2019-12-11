package runner

import (
	"github.com/jakhax/sandman/runneropt"
	"io"
)

//CodeRunner interface
type CodeRunner interface {
	Setup(opt *runneropt.Opt) (err error)
	Files(opt *runneropt.Opt)(err error)
	TestIntegration(opt *runneropt.Opt)(stdout,stderr io.Reader, err error)
	SolutionOnly(opt *runneropt.Opt)(stdout,stderr io.Reader, err error)
	SanitizeStdErr(stderr io.Reader)(sanStderr io.Reader,err error)
	SanitizeStdOut(stdout io.Reader)(sanStdout io.Reader,err error)
	TransformOutput(stdout,stderr io.Reader)(tStdout,tStderr io.Reader,err error)
}