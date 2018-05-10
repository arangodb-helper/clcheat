package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func fixFile(name string) {
	r, e := os.Open(name)
	if e != nil {
		return
	}
	c, e := ioutil.ReadAll(r)
	r.Close()
	if e == nil {
		s := string(c)
		s = strings.Replace(s, "/Zi", "/Z7", -1)
		_ = ioutil.WriteFile(name, []byte(s), 0644)
	}
}

func main() {
	out, _ := os.OpenFile("clcheat.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	fmt.Fprintln(out, "clcheat here:")
	for i := 0; i < len(os.Args); i++ {
		fmt.Fprintf(out, "%d: %v\n", i, os.Args[i])
	}
	env := os.Environ()
	fmt.Fprintf(out, "Env: %v\n", env)

	// Compute the name of the original executable:
	abs, err := os.Executable()
	if err != nil {
		fmt.Printf("Could not find absolute path for executable: %v\n", err)
		os.Exit(1)
	}
	dir := filepath.Dir(abs)
	newPath := filepath.Join(dir, "clorig.exe")

	fmt.Fprintf(out, "Actual executable to be started: %s\n", newPath)

	if os.Getenv("CLCHEAT") != "" {
		for i := 0; i < len(os.Args); i++ {
			if os.Args[i] == "/Zi" {
				os.Args[i] = "/Z7"
				fmt.Fprintf(out, "Replacing /Zi by /Z7 in argument %d\n", i)
			} else if os.Args[i][0] == '@' {
				fmt.Fprintf(out, "Replacing /Zi by /Z7 in file %s\n", os.Args[i][1:])
				fixFile(os.Args[i][1:])
			}
		}
	}

	cmd := exec.Command(newPath, os.Args[1:]...)
	cmd.Args[0] = abs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Run()
	if e != nil {
		fmt.Printf("Error in Run: %v\n", e)
	}
	fmt.Fprintln(out, "Done clcheat")
	out.Close()
}
