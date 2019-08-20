// encode or decode an IBM WebSphere Application Server password
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"
)

const maxlen = 128

var allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!()-._`~@#"
var usage = `usage: wp [OPTION] [plain text password] | [encoded password]

Options:
  -encode  Plain text password to encode. The password must be ASCII,
           no longer than 128 characters and cannot contain a space.
  -decode  WebSphere encoded password to decode. The password may be
           prefixed with {xor}.
  -h       Prints this help.

Only one option may be specified.
`

func checkPwd(s string) bool {
	var count int
	if strings.HasPrefix(s, ".") || strings.HasPrefix(s, "-") ||
		strings.HasPrefix(s, "_") {
		return false
	}
	for _, x := range s {
		count++
		if !strings.Contains(allowed, string(x)) || count > maxlen {
			return false
		}
	}
	return true
}

func main() {
	var (
		enc = flag.String("encode", "", "")
		dec = flag.String("decode", "", "")
	)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
		os.Exit(1)
	}
	flag.Parse()

	if len(os.Args[1:]) == 0 || flag.NFlag() == 2 || flag.NArg() >= 1 {
		flag.Usage()
	}
	d := make([]byte, 0, maxlen)
	switch os.Args[1] {
	case "-encode", "--encode":
		if !checkPwd(*enc) {
			flag.Usage()
		}
		for _, x := range *enc {
			d = append(d, byte(x^'_'))
		}
		pwd := base64.StdEncoding.EncodeToString(d)
		fmt.Println("{xor}" + pwd)
	case "-decode", "--decode":
		if strings.HasPrefix(*dec, "{xor}") {
			*dec = strings.TrimPrefix(*dec, "{xor}")
		}
		d, _ = base64.StdEncoding.DecodeString(*dec)
		for _, x := range d {
			fmt.Printf("%s", string(x^'_'))
		}
		fmt.Println()
	default:
		flag.Usage()
	}
}
