package main

import (
	"os/exec"
)

type SysCall struct {
	cmd    string
	output chan string
	input  chan string
	errput chan string
}

func NewSysCall(cmd string) *SysCall {
	call := new(SysCall)
	call.cmd = cmd
	call.output = make(chan string, 1)
	call.input = make(chan string, 1)
	call.errput = make(chan string, 1)
	return call
}

func (call *SysCall) Exec() error {
	cmd := exec.Command("/bin/bash", "-c", call.cmd)

	err := cmd.Run()
	if err != nil {
		return err
	}
}

func (call *SysCall) ExecTimeOut() error {

}

func (call *SysCall) AsyncExec() error {

}

func (call *SysCall) Kill() {

}
