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
  fatalErrorCheck(err)
	homeDir := user.HomeDir

	keyDir := homeDir+"/.LegalCrypt/keys"
  initFiles := []string{
    keyDir + "/id_rsa",
    keyDir + "/id_rsa.pub",
    keyDir + "/nacl"}

  testFilesDir := homeDir + "/go/src/cryptengine/test/testFiles"
	randosFiles := []string{
    testFilesDir + "/benchmark-file-megs1.rando",
    testFilesDir + "/benchmark-file-megs2.rando",
    testFilesDir + "/benchmark-file-megs15.rando",
    testFilesDir + "/benchmark-file-megs60.rando",
    testFilesDir + "/benchmark-file-megs120.rando",
    testFilesDir + "/benchmark-file-megs240.rando",
    testFilesDir + "/benchmark-file-megs512.rando",
    testFilesDir + "/benchmark-file-megs740.rando",
    testFilesDir + "/CHECKSUM.SHA512-benchmark-file"}


  fmt.Println("********* BEGIN CRYPTENGINE INTEGRATION TESTS ********")

  fmt.Println("setting up test envrionment...")
	// env init
  // TODO: remove old encrypted files

  // setting up randos
  fmt.Println("Checking randos files...")
  if fileCheck (randosFiles) != 0 {
    err = os.RemoveAll(homeDir+"/.LegalCrypt")
    fatalErrorCheck(err)
    fmt.Println("Files missing, regenerating test files (this takes awhile)...")

    _ = runCommand("./run_randos.sh")
  }

  fmt.Println("Setup complete, starting tests...")

	// keygen and setup
	fmt.Println("Generating keys...")
	output := runCommand("./cryptengine", "-gen", "-t", "rsa", "-p", password, "-eml",
		email)
	fmt.Println(string(output))

	// check keys exist
	fmt.Println("Verifying keys exists...")
  missingFiles := fileCheck (initFiles)
  fmt.Println(strconv.Itoa(missingFiles) + " files missing")

	// cleanup
	os.RemoveAll(homeDir+"/.LegalCrypt")
}

/*
takes a list of file paths and check if the files exists
prints when a file path is verified
returns # of files it failed to find
*/
func fileCheck (filepaths []string) (int){
  failCount := 0
  for _, path := range filepaths {
    _, err := os.Stat(path)
    if !os.IsNotExist(err) {
      fmt.Println("Found: " + path)
    } else if os.IsNotExist(err) {
      fmt.Println("ERROR: Unable to find " + path)
      failCount++
    } else {
      fmt.Println("ERROR: Unknown error locating " + path)
      failCount++
    }
  }
  return failCount
}

/*
wrapper for running commands via os/exec
returns output of command run
*/
func runCommand (name string, arg ...string) ([]byte){
  output, err := exec.Command(name , arg...).CombinedOutput()
  if err != nil {
    os.Stderr.WriteString(err.Error())
  }
  return output
}

/*
checks error, stops if not nil
stops test, only use for fatal errors
*/
func fatalErrorCheck(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
