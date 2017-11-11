// The following comments a pulled directly from the scrypt.Key
// source on golang's github, as of Aug 5, 2017:
// https://github.com/golang/crypto/blob/master/scrypt/scrypt.go
//
// ==============================================================================
//
// Key derives a key from the password, salt, and cost parameters, returning
// a byte slice of length keyLen that can be used as cryptographic key.
//
// N is a CPU/memory cost parameter, which must be a power of two greater than 1.
// r and p must satisfy r * p < 2³⁰. If the parameters do not satisfy the
// limits, the function returns a nil byte slice and an error.
//
// For example, you can get a derived key for e.g. AES-256 (which needs a
// 32-byte key) by doing:
//
//      dk, err := scrypt.Key([]byte("some password"), salt, 16384, 8, 1, 32)
//
// The recommended parameters for interactive logins as of 2009 are N=16384,
// r=8, p=1. They should be increased as memory latency and CPU parallelism
// increases. Remember to get a good random salt.
//
// ==============================================================================
//
// In our implementation, I'm testing with the following values:
// N := 512 * 1024  // CPU/memory cost parameter
// r := 19          // Latency of Memory Subsystem
// p := 2           // Parallelism
//
// On a Late 2013 MBP w/ 2.4GHz i5 & 8GB DDR3 @ 1600MHz, the following average
// results have been collected over a set of 10 runs / setting:
//
// N           | r   | p   | time (s)
// ============|=====|=====|=====================
// 512 * 1024  | 16  | 2   | 7.261s
// 512 * 1024  | 16  | 4   | 14.977
// 512 * 1024  | 16  | 8   | 30.542
// 512 * 1024  | 19  | 2   | 9.488
// 512 * 1024  | 19  | 4   | 19.402
//
// On the relationships between these values, and given current Apple hardware:
// maxInt = (2 * 1024 * 1024 * 1024) - 1
// r <= maxInt / 128 / p
// r <= maxInt / 256
// N <= maxInt / 128 / r
//

package main

import (
	//"encoding/base64"
	"golang.org/x/crypto/scrypt"
)

func scryptify(pass string, email string, keyLen int) []byte {
	passwd := []byte(pass)
	salt := []byte(email)
	N := 512 * 1024
	r := 19
	p := 2

	key, _ := scrypt.Key(passwd, salt, N, r, p, keyLen)

	return key
}

//func scryptify64(pass string, email string, keyLen int) string {
//	return base64.StdEncoding.EncodeToString(scryptify(pass, email, keyLen))
//}
