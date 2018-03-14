package main

import (
	"os/exec"
	"os"
	"fmt"
	"os/user"
  "strconv"
)



func main() {
	email := "test@email.com"
	password := "password"
	user, err := user.Current()
	homeDir := user.HomeDir
	keyDir := homeDir+"/.LegalCrypt/keys"
	initFiles := []string{
	  keyDir + "/id_rsa",
	  keyDir + "/id_rsa.pub",
	  keyDir + "/nacl"}

	// keygen and setup
	output, err := exec.Command("./cryptengine", "-gen", "-t", "rsa", "-p", password, "-eml",
		email).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))

	// check keys exist
  missingFiles := fileCheck (initFiles)
  fmt.Println(strconv.Itoa(missingFiles) + " files missing")

	// cleanup
	os.RemoveAll(homeDir+"/.LegalCrypt")

	// check for .LegalCrypt (ls check)
	output, err = exec.Command("ls", "-a", homeDir).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))
}

// takes a list of file paths and check if the files exists
// prints when a file path is verified
// returns # of files it failed to find
func fileCheck (filepaths []string) (int){
  failCount := 0
  for _, path := range filepaths {
    _, err := os.Stat(path)
    if !os.IsNotExist(err) {
      fmt.Println("Found: " + path)
    } else {
      fmt.Println("ERROR: Unable to find " + path)
      failCount++
    }
  }
  return failCount
}
