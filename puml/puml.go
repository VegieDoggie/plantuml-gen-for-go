package puml

import (
	"fmt"
	"github.com/VegetableDoggies/plantuml-gen-for-go/utils/arraysx"
	"github.com/VegetableDoggies/plantuml-gen-for-go/utils/filepathx"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"regexp"
	"strings"
	"sync"
)

type InterfaceSpec struct {
	Value ast.Spec
	Pkg   string
	Dir   string
}

type DepSpec struct {
	ID    string
	Value ast.Expr
	Kind  DepKind
}

type DepPackage struct {
	Dir     string
	Pkg     *doc.Package           // 包
	Scope   map[string]*ast.Object // 包域环境值
	DepsMap map[string][]*DepSpec  // 类的依赖列表
}

func (d *DepPackage) AddDepSpec(name string, depSpec *DepSpec) {
	if depSpec.Value != nil {
		d.DepsMap[name] = append(d.DepsMap[name], depSpec)
	}
}

type Portrait struct {
	Wg            *sync.WaitGroup
	Mutex         *sync.Mutex
	Rootdir       string
	Excldirs      []string
	DepPkgs       []*DepPackage
	DirPkgMap     map[string]*DepPackage
	DirsMap       map[string][]string
	InterfaceList []InterfaceSpec
	NoScope       []string
	NoBuild       []string
	Puml          string
	IsFlat        bool
}

func NewDepSpec(expr ast.Expr, isSupper bool, isWeek bool) (depSpec *DepSpec) {
	if isSupper {
		if isWeek {
			depSpec = &DepSpec{Value: expr, Kind: SupperWeek}
		} else {
			depSpec = &DepSpec{Value: expr, Kind: SupperStrong}
		}
	} else {
		depSpec = ParseDep(expr, false, isWeek)
	}
	depSpec.ID = types.ExprString(depSpec.Value)
	return
}

func NewPortrait(rootdir string, excldirs []string, isFlat bool) *Portrait {
	portrait := &Portrait{
		Wg:        new(sync.WaitGroup),
		Mutex:     new(sync.Mutex),
		Rootdir:   rootdir,
		Excldirs:  excldirs,
		DirPkgMap: make(map[string]*DepPackage),
		DirsMap:   make(map[string][]string),
		IsFlat:    isFlat,
	}
	portrait.Scan()
	portrait.DrawPuml()
	return portrait
}

func (p *Portrait) Scan() {
	p.Wg.Add(1)
	p.scan(p.Rootdir)
	p.Wg.Wait()
}

// scan - 递归扫描解析包，除了跳过排除列表，解析过程还将跳过前缀为"."的文件夹
func (p *Portrait) scan(dir string) {
	if arraysx.Index(p.Excldirs, func(excldir string) bool { return excldir == dir }) != -1 {
		p.Wg.Done()
		return
	}
	bpkg, err := build.ImportDir(dir, build.ImportComment)
	if err == nil {
		fset := token.NewFileSet()
		apkgs, _ := parser.ParseDir(fset, dir, func(info fs.FileInfo) bool {
			for _, name := range bpkg.GoFiles {
				if name == info.Name() {
					return true
				}
			}
			for _, name := range bpkg.CgoFiles {
				if name == info.Name() {
					return true
				}
			}
			return false
		}, parser.ParseComments)
		if len(apkgs) > 0 {
			docPkg := doc.New(apkgs[bpkg.Name], bpkg.ImportPath, doc.AllDecls)
			depPkg := &DepPackage{Dir: dir, Pkg: docPkg, Scope: make(map[string]*ast.Object), DepsMap: make(map[string][]*DepSpec)}
			for _, file := range apkgs[bpkg.Name].Files {
				for k := range file.Scope.Objects {
					depPkg.Scope[k] = file.Scope.Objects[k]
				}
			}
			var itfcList []InterfaceSpec
			for _, dtyp := range docPkg.Types {
				for _, spec := range dtyp.Decl.Specs {
					name := any(spec).(*ast.TypeSpec).Name.Name
					switch typ := any(spec).(*ast.TypeSpec).Type.(type) {
					case *ast.StructType:
						for _, typi := range typ.Fields.List {
							if len(typi.Names) > 0 {
								depPkg.AddDepSpec(name, NewDepSpec(typi.Type, false, false))
							} else {
								depPkg.AddDepSpec(name, NewDepSpec(typi.Type, true, false))
							}
						}
						if len(dtyp.Methods) > 0 {
							for _, method := range dtyp.Methods {
								if method.Decl.Type.Params != nil {
									for _, typi := range method.Decl.Type.Params.List {
										depPkg.AddDepSpec(name, NewDepSpec(typi.Type, false, true))
									}
								}
								if method.Decl.Type.Results != nil {
									for _, typi := range method.Decl.Type.Results.List {
										depPkg.AddDepSpec(name, NewDepSpec(typi.Type, false, true))
									}
								}
							}
						}
					case *ast.InterfaceType:
						itfcList = append(itfcList, InterfaceSpec{Value: spec, Pkg: docPkg.Name, Dir: dir})
						for _, typi := range typ.Methods.List {
							if typi.Names == nil {
								depPkg.AddDepSpec(name, NewDepSpec(typi.Type, true, true))
							}
						}
					default:
						// 语法: type Xxx abc, 类型:*ast.Ident|*ast.ArrayType..等, 如: type Kind int, type Account []string
						depPkg.AddDepSpec(name, NewDepSpec(typ, false, false))
					}
				}
			}

			p.Mutex.Lock()
			p.InterfaceList = append(p.InterfaceList, itfcList...)
			if len(depPkg.Scope) > 0 {
				p.DepPkgs = append(p.DepPkgs, depPkg)
				p.DirsMap[docPkg.Name] = append(p.DirsMap[docPkg.Name], dir)
				p.DirPkgMap[dir] = depPkg
			} else {
				p.NoScope = append(p.NoScope, dir)
			}
			p.Mutex.Unlock()
		} else {
			p.Mutex.Lock()
			p.NoBuild = append(p.NoBuild, dir)
			p.Mutex.Unlock()
		}
	} else {
		p.Mutex.Lock()
		p.NoBuild = append(p.NoBuild, dir)
		p.Mutex.Unlock()
	}
	subs := filepathx.Sub1Dirs(dir, true)
	p.Wg.Add(len(subs))
	for i := range subs {
		go p.scan(subs[i])
	}
	p.Wg.Done()
}

// PkgUniqueName - 生成包的固定唯一名，解决多个同名包问题
func (p *Portrait) PkgUniqueName(pkgName, dir string) string {
	if len(p.DirsMap[pkgName]) > 1 {
		for i := range p.DirsMap[pkgName] {
			if dir == p.DirsMap[pkgName][i] {
				return fmt.Sprintf("%v.%s", i, pkgName)
			}
		}
	}
	return pkgName
}

// ClassUniqueName - 生成含固定唯一包前缀的类名，若不指定dir将跳过包名生成逻辑
func (p *Portrait) ClassUniqueName(pkgName, dir, clss string) string {
	if dir == "" {
		return pkgName + "." + clss
	}
	return p.PkgUniqueName(pkgName, dir) + "." + clss
}

// isSelector - 判断是否为链式类型
func isSelector(depSpec *DepSpec) bool {
	_, ok := depSpec.Value.(*ast.SelectorExpr)
	return ok
}

// isSameDep - 判断两个依赖是否相同
func isSameDep(a, b *DepSpec) bool {
	return a.ID == b.ID
}

// IndentWithPreDirs - 基于当前DepPkg的索引计算缩进和父路径列表
func (p *Portrait) IndentWithPreDirs(i int, dir string) (indent0, indent1 string, preDirs []string) {
	for _, pdp := range p.DepPkgs[:i] {
		if strings.HasPrefix(dir, pdp.Dir) {
			indent0 += "\t"
			preDirs = append(preDirs, pdp.Dir)
		}
	}
	indent1 = indent0 + "\t"
	return
}

// DrawPuml - 线程保护生成puml文本
func (p *Portrait) DrawPuml() {
	p.Mutex.Lock()
	var builder strings.Builder
	builder.WriteString("@startuml\n\n")
	for i, dp := range p.DepPkgs {
		pkgUniqueName := p.PkgUniqueName(dp.Pkg.Name, p.DepPkgs[i].Dir)
		indent0, indent1, preDirs := p.IndentWithPreDirs(i, dp.Dir)
		pkgString := func() string {
			var builderPkg, builderDep strings.Builder
			if !p.IsFlat {
				builderPkg.WriteString(fmt.Sprintf("%spackage %s {\n", indent0, pkgUniqueName))
			}
			builderPkg.WriteString(fmt.Sprintf("' %s\n", dp.Dir))
			builderPkg.WriteString(DocGlobalScopeString(dp.Pkg.Consts, dp.Pkg.Vars, dp.Pkg.Funcs, indent1, pkgUniqueName))
			for _, typ := range dp.Pkg.Types {
				clssUniqueNameA := p.ClassUniqueName(pkgUniqueName, "", typ.Name)
				sFull, nsFull := arraysx.Split(dp.DepsMap[typ.Name], isSelector)
				selectors := arraysx.RemoveRedundant(sFull, isSameDep)
				for _, dep := range selectors {
					sName, sel := dep.Value.(*ast.SelectorExpr).X.(*ast.Ident).Name, types.ExprString(dep.Value.(*ast.SelectorExpr).Sel)
					for _, dir := range p.DirsMap[sName] {
						_, ok := p.DirPkgMap[dir].Scope[sel]
						if ok {
							clssUniqueNameB := p.ClassUniqueName(p.DirPkgMap[dir].Pkg.Name, dir, sel)
							if len(preDirs) > 0 {
								// 说明: if 分支对未定义的包进行空写(防止图形软件解析错误)， 如: package xxx {}
								flag := "' " + preDirs[len(preDirs)-1] + "\n"
								split := strings.Split(builder.String(), flag)
								prePkgID := fmt.Sprintf(SPkgHeader, sName)
								if strings.Index(builderPkg.String(), prePkgID) == -1 && len(split) > 0 && strings.Index(split[0], prePkgID) == -1 {
									builderPkg.WriteString(indent1 + prePkgID + "}\n")
								}
							}
							depNums := len(arraysx.Filter(sFull, func(a *DepSpec) bool { return isSameDep(a, dep) }))
							switch {
							case dep.Kind == SpecWeek && depNums > 1:
								builderDep.WriteString(DepString(clssUniqueNameA, clssUniqueNameB, indent1, SpecStrong) + "\n")
							case dep.Kind == ArrayWeek && depNums > 1:
								builderDep.WriteString(DepString(clssUniqueNameA, clssUniqueNameB, indent1, ArrayStrong) + "\n")
							default:
								builderDep.WriteString(DepString(clssUniqueNameA, clssUniqueNameB, indent1, dep.Kind) + "\n")
							}
							break
						}
					}
				}
				nonSelectors := arraysx.RemoveRedundant(nsFull, isSameDep)
				for _, dep := range nonSelectors {
					_, ok := dp.Scope[dep.ID]
					if ok && dep.ID != typ.Name {
						clssUniqueNameB := p.ClassUniqueName(dp.Pkg.Name, dp.Dir, dep.ID)
						supperI := arraysx.Index(nsFull, func(a *DepSpec) bool { return isSameDep(dep, a) && (a.Kind == SupperStrong || a.Kind == SupperWeek) })
						if supperI != -1 {
							builderDep.WriteString(DepString(clssUniqueNameA, clssUniqueNameB, indent1, nsFull[supperI].Kind) + "\n")
						} else {
							depNums := len(arraysx.Filter(nsFull, func(a *DepSpec) bool { return isSameDep(a, dep) }))
							switch {
							case dep.Kind == SpecWeek && depNums > 1:
								builderDep.WriteString(DepString(clssUniqueNameA, clssUniqueNameB, indent1, SpecStrong) + "\n")
							case dep.Kind == ArrayWeek && depNums > 1:
								builderDep.WriteString(DepString(clssUniqueNameA, clssUniqueNameB, indent1, ArrayStrong) + "\n")
							default:
								builderDep.WriteString(DepString(clssUniqueNameA, clssUniqueNameB, indent1, dep.Kind) + "\n")
							}
						}
					}
				}
				builderPkg.WriteString(builderDep.String())
				builderPkg.WriteString(DocTypeString(typ, indent1, pkgUniqueName))
				builderDep.Reset()
			}
			if !p.IsFlat {
				builderPkg.WriteString(fmt.Sprintf("%s}\n", indent0))
			}
			return builderPkg.String()
		}
		if len(preDirs) > 0 {
			fullTemp, pkgTemp := builder.String(), pkgString()
			flag := "' " + preDirs[len(preDirs)-1] + "\n"
			builder.Reset()
			builder.WriteString(strings.Replace(fullTemp, flag, pkgTemp+flag, 1))
		} else {
			builder.WriteString(pkgString())
		}
	}
	builder.WriteString("\n@enduml")
	p.Puml = regexp.MustCompile(`(?m)(^' [\S ]+\n)`).ReplaceAllString(builder.String(), "")
	p.Mutex.Unlock()
}

// ParseDep 解析依赖信息
func ParseDep(expr ast.Expr, isArrayed bool, isWeek bool) *DepSpec {
	switch typ := expr.(type) {
	default:
		// 依赖忽略: interface{} | struct{} | func(){}
		return &DepSpec{}
	case *ast.Ident, *ast.SelectorExpr:
		if isArrayed {
			if isWeek {
				return &DepSpec{Value: expr, Kind: ArrayWeek}
			} else {
				return &DepSpec{Value: expr, Kind: ArrayStrong}
			}
		} else {
			if isWeek {
				return &DepSpec{Value: expr, Kind: SpecWeek}
			} else {
				return &DepSpec{Value: expr, Kind: SpecStrong}
			}
		}
	case *ast.ParenExpr:
		return ParseDep(typ.X, isArrayed, isWeek)
	case *ast.StarExpr:
		return ParseDep(typ.X, isArrayed, isWeek)
	case *ast.Ellipsis:
		return ParseDep(typ.Elt, true, isWeek)
	case *ast.ArrayType:
		return ParseDep(typ.Elt, true, isWeek)
	case *ast.MapType:
		return ParseDep(typ.Value, true, isWeek)
	case *ast.ChanType:
		return ParseDep(typ.Value, true, isWeek)
	}
}
