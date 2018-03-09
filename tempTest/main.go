package main

import (
	"os/exec"
	"os"
	"fmt"
	"os/user"
)



func main() {
	email := "test@email.com"
	password := "password"
	user, err := user.Current()
	homeDir := user.HomeDir

	// keygen and setup
	output, err := exec.Command("./cryptengine", "-gen", "-t", "rsa", "-p", password, "-eml",
		email).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))

	// cleanup
	output, err = exec.Command("rm", "-rf", homeDir+"/.LegalCrypt").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))

	// check for .LegalCrypt (ls check)
	output, err = exec.Command("ls", "-a", homeDir).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))
}