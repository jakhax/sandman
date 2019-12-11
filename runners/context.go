package runners;

import (
	"context"
	"errors"
)

const (
	OPT="opt"
)

func CtxGetOpt(ctx context.Context) (opt *Opt, err error){
	if opt,ok := ctx.Value(OPT).(*Opt); !ok{
		err =  errors.New("Context missing Opt");
	}
	return
}