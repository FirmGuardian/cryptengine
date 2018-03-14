package main

import (
	"os/exec"
	"os"
	"fmt"
	"os/user"
  "strconv"
  "log"
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


  fmt.Println("********* BEGIN CRYPTENGINE INTEGRATION TESTS ********")

  fmt.Println("setting up test envrionment...")
	// env init
  err = os.RemoveAll(homeDir+"/.LegalCrypt")
  errorCheck(err)
  // TODO: remove old encrypted files

  // setting up randos
  // TODO: add check to see if randos needs a regen before doing it
  fmt.Println("Setting up randos files (this takes awhile)...")
  output, err := exec.Command("./run_randos.sh").CombinedOutput()
  if err != nil {
    os.Stderr.WriteString(err.Error())
  }
  fmt.Println("Setup complete, starting tests...")

	// keygen and setup
	fmt.Println("Generating keys...")
	output, err = exec.Command("./cryptengine", "-gen", "-t", "rsa", "-p", password, "-eml",
		email).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))

	// check keys exist
	fmt.Println("Verifying keys exists...")
  missingFiles := fileCheck (initFiles)
  fmt.Println(strconv.Itoa(missingFiles) + " files missing")

	// cleanup
	os.RemoveAll(homeDir+"/.LegalCrypt")

	/*
	// check for .LegalCrypt (ls check)
	output, err = exec.Command("ls", "-a", homeDir).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))*/
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

// checks error, stops if not nil
func errorCheck(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
