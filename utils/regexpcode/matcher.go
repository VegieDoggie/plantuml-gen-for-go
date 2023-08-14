package regexpcode

import (
	"bytes"
	"fmt"
	"github.com/VegetableDoggies/plantuml-gen-for-go/utils/arraysx"
	"go/format"
	"log"
	"strings"
)

var (
	space  = []byte(" ")
	braceR = []byte(")")
	eq     = []byte("=")
	iota   = []byte("iota")
)

type Struct struct {
	Name    string
	Package string
	Supers  []string
	Fields  []string
	Funcs   []*MFunc
}

type GFunc struct {
	Name   string
	Header string
	Types  []string
}

type MFunc struct {
	*GFunc
	Master string
}

type Interface struct {
	Name   string
	Supers []string
	Funcs  []*GFunc
}

func ClearLineComment(bs []byte) []byte {
	return LineComment.ReplaceAll(bs, Null)
}

func ClearMultiComment(bs []byte) []byte {
	return MultiComment.ReplaceAll(bs, Null)
}

func ClearEmptyLine(bs []byte) []byte {
	return EmptyLineCompile.ReplaceAll(bs, Null)
}

func ClearAnnotation(bs []byte) []byte {
	return AnnotationCompile.ReplaceAll(bs, Null)
}

func ClenCode(bs []byte) []byte {
	// 代码风格格式化
	source, err := format.Source(ClearAnnotation(ClearLineComment(ClearMultiComment(bs))))
	if err != nil {
		log.Fatal("lexer.Lex::", err)
	}
	return ClearEmptyLine(source)
}

func FindPackage(content []byte) string {
	return string(PackageCmp.FindSubmatch(content)[1])
}

func FindAllImports(content []byte) (imports []string) {
	blocks := ImportBlockCmp.FindAll(content, -1)
	for _, item := range blocks {
		subs := ImportBlockSubCmp.FindAllSubmatch(item, -1)
		for i := 0; i < len(subs); i++ {
			imports = append(imports, string(subs[i][1]))
		}
	}
	return imports
}

func FindAllVars(content []byte) (vars []string) {
	subs := VarLineCmp.FindAllSubmatch(content, -1)
	for i := 0; i < len(subs); i++ {
		vars = append(vars, string(subs[i][1]))
	}
	blocks := VarBlockCmp.FindAll(content, -1)
	for _, item := range blocks {
		subs = VarBlockLineCmp.FindAllSubmatch(item, -1)
		for i := 0; i < len(subs); i++ {
			vars = append(vars, string(subs[i][1]))
		}
	}
	return vars
}

func FindAllConsts(content []byte) (consts []string) {
	subs := ConstLineCmp.FindAllSubmatch(content, -1)
	for i := 0; i < len(subs); i++ {
		consts = append(consts, string(subs[i][1]))
	}
	blocks := ConstBlockCmp.FindAll(content, -1)
	for _, item := range blocks {
		subs = ConstBlockLineCmp.FindAllSubmatch(item, -1)
		for i, isIota := 0, false; i < len(subs); i++ {
			if isIota && bytes.Index(subs[i][1], eq) == -1 {
				consts = append(consts, string(subs[i][1])+" [iota]")
				continue
			}
			if bytes.HasSuffix(subs[i][1], iota) {
				isIota = true
			}
			consts = append(consts, string(subs[i][1]))
		}
	}
	return consts
}

func FindAllStructs(content []byte, pkg string) (structs []*Struct) {
	// 1- 解析单行类型
	subs := TypeCmp.FindAllSubmatch(content, -1)
	for i := 0; i < len(subs); i++ {
		structs = append(structs, &Struct{Name: string(subs[i][1]), Supers: []string{string(subs[i][2])}})
	}
	// 2- 解析块类型
	blocks := StructCmp.FindAllSubmatch(content, -1)
	for i := 0; i < len(blocks); i++ {
		sts := &Struct{Name: string(blocks[i][1]), Package: pkg}
		supers := TypeSuperCmp.FindAllSubmatch(blocks[i][2], -1)
		for j := 0; j < len(supers); j++ {
			sts.Supers = append(sts.Supers, string(supers[j][1]))
		}
		fields := StructFieldCmp.FindAllSubmatch(blocks[i][2], -1)
		for j := 0; j < len(fields); j++ {
			sts.Fields = append(sts.Fields, string(fields[j][1])+" "+string(fields[j][2]))
		}
		structs = append(structs, sts)
	}
	return structs
}

func FindAllMFuncs(content []byte) (mFuncs []*MFunc) {
	lines := StructFuncCmp.FindAllSubmatch(content, -1)
	for i := 0; i < len(lines); i++ {
		var fc *MFunc
		if bytes.Index(lines[i][1], space) == -1 {
			fc.Master = strings.ReplaceAll(string(lines[i][1]), "*", "")
		} else {
			fc.Master = strings.ReplaceAll(string(bytes.Split(lines[i][1], space)[1]), "*", "")
		}
		fc.Name, fc.Header = string(lines[i][2]), string(lines[i][2])+string(lines[i][3])
		split := bytes.Split(lines[i][3], braceR)
		fc.Types = findAllFuncDeps(string(split[0]), string(split[1]))
		mFuncs = append(mFuncs, fc)
	}
	return mFuncs
}

// BindingStructs 检索并绑定成员函数，每次匹配都会移除@mFuncsPointer列表中的元素
func BindingStructs(structs []*Struct, mFuncsPointer *[]*MFunc) {
	mFuncs := *mFuncsPointer
	for i := 0; i < len(mFuncs); {
		for j := 0; j < len(structs); j++ {
			if mFuncs[i].Master == structs[j].Name {
				structs[j].Funcs = append(structs[j].Funcs, mFuncs[i])
				mFuncs = arraysx.RemoveByIndex(mFuncs, i)
				goto NextFc
			}
		}
		i++
	NextFc:
	}
	*mFuncsPointer = mFuncs
}

func FindAllInterfaces(content []byte) (interfaces []*Interface) {
	blocks1 := InterfaceCmp.FindAllSubmatch(content, -1)
	for i := 0; i < len(blocks1); i++ {
		itc := &Interface{Name: string(blocks1[i][1])}
		// 1- 解析父类
		supers := TypeSuperCmp.FindAllSubmatch(blocks1[i][2], -1)
		for j := 0; j < len(supers); j++ {
			itc.Supers = append(itc.Supers, string(supers[j][1]))
		}
		// 2- 逐行解析函数和依赖
		lines := InterfaceLineCmp.FindAllSubmatch(blocks1[i][2], -1)
		for j := 0; j < len(lines); j++ {
			itc.Funcs = append(itc.Funcs, &GFunc{
				Name:   string(lines[j][1]),
				Header: strings.TrimSpace(string(lines[j][0])),
				Types:  findAllFuncDeps(string(lines[j][2]), string(lines[j][3]))})
		}
		interfaces = append(interfaces, itc)
	}
	return interfaces
}

func FindAllGFuncs(content []byte) (v []*GFunc) {
	fcs := FuncGlobalCmp.FindAllSubmatch(content, -1)
	for i := 0; i < len(fcs); i++ {
		fmt.Println(i, fcs[i][1], fcs[i][2], fcs[i][0])
		split := strings.Split(string(fcs[i][2]), ")")
		v = append(v, &GFunc{
			Name:   string(fcs[i][1]),
			Header: string(fcs[i][1]) + string(fcs[i][2]),
			Types:  findAllFuncDeps(split[0], split[1])})
	}
	return v
}

func findAllFuncDeps(params string, returns string) (v []string) {
	types := InterfaceDepsCmp.FindAllStringSubmatch(params, -1)
	for i := 0; i < len(types); i++ {
		v = append(v, types[i][1])
	}
	returns = strings.TrimSpace(returns)
	if len(returns) > 0 {
		if strings.Index(returns, "(") == -1 {
			v = append(v, returns)
		} else {
			if InterfaceDepsCmp.MatchString(returns) {
				types = InterfaceDepsCmp.FindAllStringSubmatch(returns, -1)
				for i := 0; i < len(types); i++ {
					v = append(v, types[i][1])
				}
			} else {
				types := strings.Split(strings.ReplaceAll(strings.ReplaceAll(returns, "(", ""), ")", ""), ", ")
				for i := 0; i < len(types); i++ {
					v = append(v, types[i])
				}
			}
		}
	}
	return v
}
