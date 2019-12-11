package sandbox

import (
	"github.com/jakhax/sandman/runneropt"
)

// MissingGitURL error
type MissingGitURL struct{
	
}
func (e MissingGitURL) Error() string{
	return "Missing Git url"
}

// DownloadFromGit get project from git
func DownloadFromGit(opt *runneropt.Opt ,timeout int)(err error){

	spwanOpt := &SpawnOpt{
		Dir:opt.Dir,
		Env:opt.Env,
	}
	if timeout != 0{
		spwanOpt.Timeout = timeout
	}
	url :=  opt.GitURL;
	if url == ""{
		err =  MissingGitURL{}
		return;
	}
	args := []string{"clone","--depth=1",url,opt.Dir+"/."}
	_,_,err = Spwan(spwanOpt,"git",args,nil)
	return
}