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

  // key locations
  keyDir := homeDir + "/.LegalCrypt/keys"
  initFiles := []string{
    keyDir + "/id_rsa",
    keyDir + "/id_rsa.pub",
    keyDir + "/nacl"}

  // test file locations
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

  // Documents directories
  documentsDir := homeDir + "/Documents/LegalCrypt"
  securedDir := documentsDir + "/Secured"
  //receivedDir := documentsDir + "/Received"


  fmt.Println("********* BEGIN CRYPTENGINE INTEGRATION TESTS ********")

  fmt.Println("setting up test envrionment...")
  err = os.RemoveAll(homeDir+"/.LegalCrypt")
  err = os.RemoveAll(documentsDir)
  fatalErrorCheck(err)

  // TODO: remove old encrypted files

  // setting up randos
  fmt.Println("Checking randos files...")
  if fileCheck (randosFiles) != 0 {
    os.RemoveAll(testFilesDir)
    fmt.Println("Files missing, regenerating test files (this takes awhile)...")
    _, err = runCommand("./run_randos.sh")
    fatalErrorCheck(err)
  }

  fmt.Println("Setup complete, starting tests...")

	// keygen and setup
	fmt.Println("Generating keys...")
	output, err := runCommand("./cryptengine", "-gen", "-t", "rsa", "-p", password, "-eml",
		email)
	fmt.Println(string(output))

	// check keys exist
	fmt.Println("Verifying keys exists...")
  missingFiles := fileCheck (initFiles)
  fmt.Println(strconv.Itoa(missingFiles) + " files missing")

  // TODO: Update test to use pre-existing file, keys folder, and checksum for encrypted version of file
  // ^ might not work due to salt
  fmt.Println("Testing Encryption...")
  oneMegFile :=  "benchmark-file-megs1.rando"
  oneMegEncrypted := oneMegFile+".lcsf"
  output, _ = runCommand("./cryptengine", "-e", "-t", "rsa", testFilesDir+"/"+oneMegFile)
  fmt.Println(string(output))
  missingFiles = fileCheck([]string{securedDir+"/"+oneMegEncrypted})
  if missingFiles != 0 {
    fmt.Println("ERROR: Unable to find " + oneMegEncrypted)
  } else {
    fmt.Println("SUCCESS: " + oneMegEncrypted + " found!")
  }

  // TODO: TEST DECRYPT: encrypt and decrypt file, verify checksum
  // TODO: Replace w/ stored .lcsf file and keys folder for test so there's no reliance on encrypt functionality
  fmt.Println("Testing Decryption...")
  twoMegFile := "benchmark-file-megs2.rando"
  //twoMegEncrypted := twoMegFile+".lcsf"
  fmt.Println("Encrypting...")
  output, _ = runCommand("./cryptengine", "-e", "-t", "rsa", testFilesDir+"/"+twoMegFile)
  fmt.Println(string(output))
  fmt.Println("Decrypting...")



  // TODO: TEST 512MB ENCRYPT: Test encrypting a 512MB file, this should succeed

  // TODO: TEST OVER LIMIT ENCRYPT: Test encrypting a file above 512MB, this should fail
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
func runCommand (name string, arg ...string) ([]byte, error){
  output, err := exec.Command(name , arg...).CombinedOutput()
  if err != nil {
    os.Stderr.WriteString(err.Error())
  }
  return output, err
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

/*
Checks checksum of given file in a given directory.  Comparision is made w/
the CHECKSUM.SHA512-benchmark-file in the testFiles folder.
*/
func verifyChecksum (filename string, fileDir string) (bool) {
  return false
}
