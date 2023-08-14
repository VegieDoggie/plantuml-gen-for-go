package puml

import (
	"fmt"
	"testing"
)

func TestNewPortrait(t *testing.T) {
	portrait := NewPortrait("D:\\Programs\\go-puml-gen", []string{"D:\\Programs\\go-puml-gen\\mod"})
	fmt.Println(portrait)
}
