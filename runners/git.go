package runners

import (
	"time"
	"context"
)

// MissingGitURL error
type MissingGitURL struct{
	
}
func (e MissingGitURL) Error() string{
	return "Missing Git url"
}

// DownloadFromGit get project from git
func DownloadFromGit(ctx context.Context,timeout time.Duration)(err error){
	opt,err :=  CtxGetOpt(ctx)
	if err != nil{
		return
	}
	spwanOpt := &SpawnOpt{
		Dir:opt.Dir,
		Env:opt.Env,
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