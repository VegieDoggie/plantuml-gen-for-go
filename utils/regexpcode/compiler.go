package regexpcode

import "regexp"

var (
	Null              = []byte("")
	LineComment       = regexp.MustCompile(`//.*`)
	MultiComment      = regexp.MustCompile(`/\*(\n|.)*?\*/`)
	EmptyLineCompile  = regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	AnnotationCompile = regexp.MustCompile("`(\\w+:\"\\w+\"\\s*)`")
	PackageCmp        = regexp.MustCompile(`(?m)^package (\w+)`)
	ImportBlockCmp    = regexp.MustCompile(`(?m)^import (\((\s+(\w+ )?("([\w/.]+)"))+\s*\)|"([\w/.]+)")`)
	ImportBlockSubCmp = regexp.MustCompile(`"([\w/.]+)"`)
	VarLineCmp        = regexp.MustCompile(`(?m)^var (\w+(?: (?:` + paramType + `)?(?: = .+?)?){?\s+$`)
	VarBlockCmp       = regexp.MustCompile(`(?m)^var \((\s+\w+([\t ]+(` + paramType + `)?([\t ]+=[\t ]+.+)?)+\s*\)`)
	VarBlockLineCmp   = regexp.MustCompile(`[\t ](\w+(?:[\t ]+(?:` + paramType + `)?(?:[\t ]+=[\t ]+[^\n\r]+)?)`)
	ConstLineCmp      = regexp.MustCompile(`(?m)^const (\w+(?: (?:` + paramType + `)?(?: = .+?)?)\s+$`)
	ConstBlockCmp     = regexp.MustCompile(`(?m)^const \((\s+\w+([\t ]+(` + paramType + `)?([\t ]+=[\t ]+.+?)?)+\s*\)`)
	ConstBlockLineCmp = VarBlockLineCmp
	TypeCmp           = regexp.MustCompile(`(?m)^type (\w+) ((?:` + paramType + `)\s+$`)
	TypeSuperCmp      = regexp.MustCompile(`(?m)^[\t ]+(\w+)\s*$`)
	StructCmp         = regexp.MustCompile(`(?m)^(?:type |[\t ]+)(\w+) struct\s?\{((?:\s+\w+(?:, \w+)*?(?:[\t ]+(?:` + paramType + `)?)+)\s*}`)
	StructFieldCmp    = regexp.MustCompile(`(?m)^[\t ]+(\w+(?:, \w+)*)[\t ]+((?:(?:` + paramType + `)+)\s*$`)
	StructFuncCmp     = regexp.MustCompile(`(?m)^func \((\*?\w+? ?\*?\w+?)\) (\w+)(\((?:\w+(?:, \w+)*? (?:chan |[\w\[\].*]|\{})+(?:, )?)*\)[\t ]*\(?(?:\w*?(?:, \w+)*? ?(?:chan |[\w\[\].*]|\{})*(?:, )?)*\)?)\s*?\{`)
	InterfaceCmp      = regexp.MustCompile(`(?m)^(?:type |[\t ]+)(\w+) interface\s?\{((?:\s*\w*(?:` + funcParamList + `)?)*)\s*}`)
	InterfaceLineCmp  = regexp.MustCompile(`[\t ]+(\w+)(\(.*?\))(.*)?`)
	InterfaceDepsCmp  = regexp.MustCompile(`\w+(?:, \w+)*? ((?:` + paramType + `)(?:, )?`)
	FuncGlobalCmp     = regexp.MustCompile(`(?m)^func (\w+)(` + funcParamList + `)\s*?\{`)
)

const (
	paramType     = `chan |[\w\[\].*]|\{})+`
	funcParamList = `\((?:\w+(?:, \w+)*? (?:` + paramType + `(?:, )?)*\)[\t ]*\(?(?:\w*?(?:, \w+)*? ?(?:chan |[\w\[\].*]|\{})*(?:, )?)*\)?`
)
