package test

import (
	"testing"

	// "." allows use of package functions w/o their package prefix
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// helper functions
func VerifyKeyGen() {
  
}

func TestCryptengine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cryptengine Integration Tests")
}

var _ = Describe("Cryptengine", func() {
  Context("generating keys", func() {
    // rm -rf ~/.LegalCrypt
    // ./cryptengine -gen -t rsa -p $password -eml $email
    It("creates an rsa keypair and salt", func() {
      // check for existence of keypair & salt
      // NOTE: this also verifies .LegalCrypt is made since keys/salt are inside
    })
  })
})
