package main
// file to test using command line from go

import (
  "fmt"
  "os"
  "os/exec"
)

const HOME string = "/Users/kevincisler"
const GOPATH string = "/Users/kevincisler/go"
const LC_CONFIG_DIR = HOME + "/.LegalCrypt"
const KEY_DIR string = LC_CONFIG_DIR + "/keys"
const LC_DOCS_DIR = HOME + "/Documents/LegalCrypt"
const SECURED_DIR string = LC_DOCS_DIR + "/Secured"
const RECEIVED_DIR string = LC_DOCS_DIR + "/Received"
const REPO_DIR string = GOPATH + "/src/cryptengine"
const TEST_DIR = REPO_DIR + "/test"

const PASSWORD = "Pass"
const EMAIL = "test@test.com"

func deleteDir (dir string) {
  if err := exec.Command("rm", "-rf", dir).Run(); err != nil {
    fmt.Fprintln(os.Stderr, err)
  }
}

func resetTest() {
  deleteDir(LC_CONFIG_DIR)
  deleteDir(LC_DOCS_DIR)
}

func main () {

  //remove old keys/direcotries
  resetTest()
}

