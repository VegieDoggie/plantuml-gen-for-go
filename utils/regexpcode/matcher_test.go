package regexpcode

import (
	"fmt"
	"testing"
)

func TestClearMultiComments(t *testing.T) {
	s := `
		start_code();
		/* First comment */
		more_code();
		/* Second

comment */
		end_code();
	`
	fmt.Println(string(ClearMultiComment([]byte(s))))
}

func TestClearLineComment(t *testing.T) {
	s := `
		start_code();
		// First comment
		more_code();// append comment
		// Second comment
		end_code();
	`
	fmt.Println(string(ClearLineComment([]byte(s))))
}

func TestClearEmptyLine(t *testing.T) {
	s := `
		start_code();

		more_code();

		end_code();
	`
	fmt.Println(string(ClearEmptyLine([]byte(s))))
}

func TestClearAnnotation(t *testing.T) {
	s := "type CryptoJSON struct {\n\tCipher       string                 `json:\"cipher\"`\n\tCipherText   string                 `json:\"ciphertext\"`\n\tCipherParams cipherparamsJSON       `json:\"cipherparams\"`\n\tKDF          string                 `json:\"kdf\"`\n\tKDFParams    map[string]interface{} `json:\"kdfparams\"`\n\tMAC          string                 `json:\"mac\"`\n}"
	fmt.Println(string(ClearAnnotation([]byte(s))))
}
