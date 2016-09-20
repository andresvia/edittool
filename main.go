package main

import (
	"bytes"
	"github.com/andresvia/editlib/editlib"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
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
	}
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) (err error) {
	edit := c.String("edit")
	ensure := c.String("ensure")

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

	if err = editlib.Edit(&out_buf, edit_buf, "# EDITTOL GENERATED DO NOT EDIT", "# EDITTOL GENERATED DO NOT EDIT", ensure_string); err != nil {
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
	err = cli.NewExitError("Edited", 2)
	return
}
