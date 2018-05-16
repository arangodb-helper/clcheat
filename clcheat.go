package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func fixFile(out *os.File, name string) {
	r, e := os.Open(name)
	if e != nil {
		return
	}
	c, e := ioutil.ReadAll(r)
	r.Close()
	if e == nil {
		s := string(c)
		s = strings.Replace(s, "/\x00Z\x00i\x00", "/\x00Z\x007\x00", -1)
		e = ioutil.WriteFile(name, []byte(s), 0664)
		if e != nil {
			fmt.Fprintf(out, "Error in WriteFile: %v\n", e)
		}
		name2 := name + ".guk"
		e = ioutil.WriteFile(name2, []byte(s), 0664)
		if e != nil {
			fmt.Fprintf(out, "Error in WriteFile2: %v\n", e)
		}
		fmt.Fprintf(out, "New content of file: %s\n", s)
	}
}

func main() {
	out, e := os.OpenFile("C:/Users/Max/clcheat.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		out = os.Stdout
	}
	fmt.Fprintf(out, "clcheat here %s:", time.Now().Format(time.UnixDate))
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

	if os.Getenv("CLCHEAT") == "1" {
		for i := 0; i < len(os.Args); i++ {
			if os.Args[i] == "/Zi" {
				os.Args[i] = "/Z7"
				fmt.Fprintf(out, "Replacing /Zi by /Z7 in argument %d\n", i)
			} else if os.Args[i][0] == '@' {
				fmt.Fprintf(out, "Replacing /Zi by /Z7 in file %s\n", os.Args[i][1:])
				fixFile(out, os.Args[i][1:])
			}
		}
	}

	cmd := exec.Command(newPath, os.Args[1:]...)
	cmd.Args[0] = abs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e = cmd.Run()
	if e != nil {
		fmt.Printf("Error in Run: %v\n", e)
	}
	fmt.Fprintln(out, "Done clcheat")
	out.Close()
}
