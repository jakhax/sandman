package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/jakhax/sandman/runneropt"
	"github.com/jakhax/sandman/sandbox"
)

type Output struct{
	Stderr string `json:"stderr"`
	Stdout string `json:"stdout"`
}

func HandleJob(w http.ResponseWriter, r *http.Request){
	payload :=  &runneropt.Opt{};
	body, err:= ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(body,payload)
	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}
	err = payload.OK()
	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}
	
	
	
	s,err := sandbox.NewSandBox()
	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}
	stdout,stderr, err := s.Run(payload)
	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}
	// sandbox.WriteToStd(stdout,stderr)
	res := &Output{}
	if stdout != nil{
		stdoutB , errX := ioutil.ReadAll(stdout)
		if errX != nil {
			w.Write([]byte(errX.Error()))
			return
		} 
		fmt.Println(string(stdoutB))
		res.Stdout = string(stdoutB)
	}
	if stderr != nil{
		stderrB , errX := ioutil.ReadAll(stderr)
		if errX != nil {
			w.Write([]byte(errX.Error()))
			return
		} 
		res.Stderr = string(stderrB)
	}
	res.Stdout = sandbox.CleanTags(res.Stdout)
	fmt.Println(res.Stdout)

	resC , err :=  json.Marshal(res)
	if err != nil{
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(resC)
}

var testServerCommand = &cobra.Command{
	Use:"test_server",
	Run:func(cmd *cobra.Command, args []string){
		server := http.NewServeMux()
		server.HandleFunc("/test",HandleJob)
		err := http.ListenAndServe(":8000",server)
		if err != nil{
			fmt.Println(err.Error());
		}
		_ =1 
		return
	},
}

func init(){
	rootCmd.AddCommand(testServerCommand);
}