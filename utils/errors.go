package utils;

import(
	"log"
	"runtime"
)

func HandleError(err error){
	if(err != nil){
		_, fn, line, _ := runtime.Caller(1)
        log.Fatalf("[error] %s:%d %v", fn, line, err)
	}
}
