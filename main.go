package main

import (
	"bufio"
	"bytes"
	"github.com/andresvia/editlib/editlib"
	"gopkg.in/urfave/cli.v1"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "edit",
			Usage: "the file to edit",
		},
		cli.StringFlag{
			Name:  "ensure",
			Usage: "the file to ensure is included",
		},
		cli.StringFlag{
			Name:  "reload",
			Usage: "the command and arguments to run if there's change",
		},
	}
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) (err error) {
	edit := c.String("edit")
	ensure := c.String("ensure")
	reload := c.String("reload")

	edit_bytes := []byte{}
	if edit_bytes, err = ioutil.ReadFile(edit); err != nil {
		err = cli.NewExitError(err.Error(), 66)
		return
	}

	edit_buf := bytes.NewBuffer(edit_bytes)

	ensure_bytes := []byte{}
	if ensure_bytes, err = ioutil.ReadFile(ensure); err != nil {
		err = cli.NewExitError(err.Error(), 66)
		return
	}

	ensure_string := string(ensure_bytes)
	out_buf := bytes.Buffer{}

	if err = editlib.Edit(&out_buf, edit_buf, "# EDITTOOL GENERATED DO NOT EDIT", "# EDITTOOL GENERATED DO NOT EDIT", ensure_string); err != nil {
		err = cli.NewExitError(err.Error(), 66)
		return
	}

	info, _ := os.Stat(edit)

	out_bytes := out_buf.Bytes()
	if bytes.Equal(edit_bytes, out_bytes) {
		return
	}

	if err = ioutil.WriteFile(edit, out_bytes, info.Mode()); err != nil {
		err = cli.NewExitError(err.Error(), 66)
		return
	}

	reload_scanner := bufio.NewScanner(strings.NewReader(reload))
	reload_scanner.Split(bufio.ScanWords)

	reload_tokens := []string{}
	for reload_scanner.Scan() {
		reload_tokens = append(reload_tokens, reload_scanner.Text())
	}

	reload_args := []string{}
	if len(reload_tokens) > 1 {
		reload_args = reload_tokens[1:]
	}

	if len(reload_tokens) == 0 {
		err = cli.NewExitError("Edited", 2)
		return
	}

	cmd := exec.Command(reload_tokens[0], reload_args...)
	var stdOut io.ReadCloser
	var stdErr io.ReadCloser
	if stdOut, err = cmd.StdoutPipe(); err != nil {
		err = cli.NewExitError(err.Error(), 66)
		return
	}
	if stdErr, err = cmd.StderrPipe(); err != nil {
		err = cli.NewExitError(err.Error(), 66)
		return
	}
	if err = cmd.Start(); err != nil {
		err = cli.NewExitError(err.Error(), 66)
		return
	}

	stdOutBytes := []byte{}
	stdErrBytes := []byte{}
	go func() {
		stdOutBytes, _ = ioutil.ReadAll(stdOut)
	}()
	go func() {
		stdErrBytes, _ = ioutil.ReadAll(stdErr)
	}()

	errStr := ""
	status := 0
	if err = cmd.Wait(); err != nil {
		stdOutStr := strings.TrimSpace(string(stdOutBytes[:]))
		stdErrStr := strings.TrimSpace(string(stdErrBytes[:]))
		outArr := []string{"Reload failed"}
		if stdOutStr != "" {
			outArr = append(outArr, stdOutStr)
		}
		if stdErrStr != "" {
			outArr = append(outArr, stdErrStr)
		}
		outArr = append(outArr, err.Error())
		errStr = strings.Join(outArr, "\n")
		status = err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus()
	}
	err = cli.NewExitError(errStr, status)
	return
}
