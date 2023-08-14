package filepathx

import (
	"fmt"
	"regexp"
	"testing"
)

func TestShortWalkedPath(t *testing.T) {
	x := `
	}
' C:\Users\Administrator\Desktop\go-ethereum-master
	class ethereum.* << (G,DarkSeaGreen) >> {
		..var..
		{field} + NotFound : errors.New("not found")
	}`
	reg := regexp.MustCompile(`(?m)(^' [\S ]+\n)`)
	fmt.Println(reg.ReplaceAllString(x, ""))
}
