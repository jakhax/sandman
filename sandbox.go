package main;
import (
	"io"
	"strings"
	"fmt"
)


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