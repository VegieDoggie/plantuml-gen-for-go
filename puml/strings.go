package puml

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"go/types"
	"strings"
)

const (
	Main              = "main"
	ClassModifier     = "class"
	InterfaceModifier = "interface"
	PublicModifier    = "+"
	PrivateModifier   = "-"
	StaticModifier    = "{static}"
	FieldModifier     = "{field}"
	MethodModifier    = "{method}"
	PrototypeArea     = "..prototype.."
	ConstArea         = "..const.."
	VarArea           = "..var.."
	FuncArea          = "..func.."
	FieldArea         = "..field.."
	MethodArea        = "..method.."
	SupperArea        = "..supper.."
	ConstructorArea   = "..constructor.."
	ScopeArea         = "==scope=="
	SPkgHeader        = "package %s {"
)

type DepKind int

const (
	SupperWeek DepKind = 1 << iota
	SupperStrong
	SpecWeek
	SpecStrong
	ArrayWeek
	ArrayStrong
)

func DepString(a, b, indent string, kind DepKind) string {
	switch kind {
	case SupperWeek:
		return indent + a + " ..|> " + b
	case SupperStrong:
		return indent + a + " --|> " + b
	case SpecWeek:
		return indent + a + " ..> " + b
	case SpecStrong:
		return indent + a + " --> " + b
	case ArrayWeek:
		return indent + a + " *..> " + b
	case ArrayStrong:
		return indent + a + " *--> " + b
	default:
		return indent + a + " .. " + b
	}
}

func (k DepKind) IsStrong() bool {
	return k == SupperStrong || k == SpecStrong || k == ArrayStrong
}

func EscapeString(text string) string {
	if len(text) > 0 {
		if len(text) > 32 {
			text = text[:32] + " .."
		}
		text = strings.ReplaceAll(text, "\\r", "\\\\r")
		text = strings.ReplaceAll(text, "\\n", "\\\\n")
		text = strings.ReplaceAll(text, "\\t", "\\\\t")
		text = strings.ReplaceAll(text, "\r", "\\\\r")
		text = strings.ReplaceAll(text, "\n", "\\\\n")
		text = strings.ReplaceAll(text, "\t", "\\\\t")
	}
	return text
}

func AstValueSpecString(spec ast.Spec, sep string) string {
	if spec != nil {
		var names, values []string
		for _, ident := range spec.(*ast.ValueSpec).Names {
			names = append(names, ident.Name)
		}
		for _, expr := range spec.(*ast.ValueSpec).Values {
			values = append(values, types.ExprString(expr))
		}
		n, v := strings.Join(names, ", "), EscapeString(strings.Join(values, ", "))
		if len(v) > 0 {
			if token.IsExported(n) {
				return fmt.Sprintf("%s %s %s %s %s", FieldModifier, PublicModifier, n, sep, v)
			} else {
				return fmt.Sprintf("%s %s %s %s %s", FieldModifier, PrivateModifier, n, sep, v)
			}
		} else {
			if token.IsExported(n) {
				return fmt.Sprintf("%s %s %s", FieldModifier, PublicModifier, n)
			} else {
				return fmt.Sprintf("%s %s %s", FieldModifier, PrivateModifier, n)
			}
		}
	}
	return ""
}

func AstFuncString(fc *ast.FuncType, pkgName, fcName string) string {
	if fc != nil {
		decl := strings.Replace(types.ExprString(fc), "func", fcName, 1)
		if pkgName == Main && fcName == Main {
			return fmt.Sprintf("%s %s **%s**", MethodModifier, PublicModifier, decl)
		} else {
			if token.IsExported(fcName) {
				return fmt.Sprintf("%s %s %s", MethodModifier, PublicModifier, decl)
			} else {
				return fmt.Sprintf("%s %s %s", MethodModifier, PrivateModifier, decl)
			}
		}
	}
	return ""
}

func DocFuncString(fc *doc.Func, pkgName string) string {
	if fc != nil && fc.Decl != nil {
		return AstFuncString(fc.Decl.Type, pkgName, fc.Name)
	}
	return ""
}

func DocValuesStrings(values []*doc.Value, indent, sep string) string {
	var builder strings.Builder
	for _, value := range values {
		for _, spec := range value.Decl.Specs {
			builder.WriteString(fmt.Sprintf("%s%s\n", indent, AstValueSpecString(spec, sep)))
		}
	}
	return builder.String()
}
func DocFuncsStrings(fcs []*doc.Func, indent, pkgName string) string {
	var builder strings.Builder
	for _, fc := range fcs {
		builder.WriteString(fmt.Sprintf("%s%s\n", indent, DocFuncString(fc, pkgName)))
	}
	return builder.String()
}

func StructFieldsStrings(styp *ast.StructType, indent string) (supers, fields []string) {
	if styp != nil && styp.Fields != nil {
		for _, field := range styp.Fields.List {
			var names []string
			for _, ident := range field.Names {
				names = append(names, ident.Name)
			}
			t := types.ExprString(field.Type)
			if len(field.Names) == 0 {
				supers = append(supers, fmt.Sprintf("%s%s %s\n", indent, StaticModifier, t))
			} else {
				n := strings.Join(names, ", ")
				if token.IsExported(n) {
					fields = append(fields, fmt.Sprintf("%s%s %s %s : %s\n", indent, FieldModifier, PublicModifier, n, t))
				} else {
					fields = append(fields, fmt.Sprintf("%s%s %s %s : %s\n", indent, FieldModifier, PrivateModifier, n, t))
				}
			}
		}
	}
	return
}

func InterfaceFieldsStrings(ityp *ast.InterfaceType, indent string) (supers, methods []string) {
	if ityp != nil && ityp.Methods != nil {
		for _, method := range ityp.Methods.List {
			if len(method.Names) == 0 {
				supers = append(supers, fmt.Sprintf("%s%s %s\n", indent, StaticModifier, types.ExprString(method.Type)))
			} else {
				methods = append(methods, fmt.Sprintf("%s%s\n", indent, AstFuncString(method.Type.(*ast.FuncType), "", method.Names[0].Name)))
			}
		}
	}
	return
}

func DocGlobalScopeString(cons []*doc.Value, vrs []*doc.Value, fcs []*doc.Func, indentHead, PkgNumberName string) string {
	var builder strings.Builder
	indentBody := indentHead + "\t"
	if len(cons)+len(vrs)+len(fcs) > 0 {
		builder.WriteString(fmt.Sprintf("%s%s %s.* << (G,DarkSeaGreen) >> {\n", indentHead, ClassModifier, PkgNumberName))
		if len(cons) > 0 {
			builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, ConstArea, DocValuesStrings(cons, indentBody, "=")))
		}
		if len(vrs) > 0 {
			builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, VarArea, DocValuesStrings(vrs, indentBody, ":")))
		}
		if len(fcs) > 0 {
			builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, FuncArea, DocFuncsStrings(fcs, indentBody, SelectorTail(PkgNumberName))))
		}
		builder.WriteString(fmt.Sprintf("%s}\n", indentHead))
	}
	return builder.String()
}

func SelectorTail(PkgNumberName string) string {
	if strings.Index(PkgNumberName, ".") != -1 {
		split := strings.Split(PkgNumberName, ".")
		return split[len(split)-1]
	}
	return PkgNumberName
}

func DocTypeString(dtyp *doc.Type, indentHead string, PkgNumberName string) string {
	var builder strings.Builder
	indentBody := indentHead + "\t"
	for _, spec := range dtyp.Decl.Specs {
		name := any(spec).(*ast.TypeSpec).Name.Name
		switch typ := any(spec).(*ast.TypeSpec).Type.(type) {
		case *ast.InterfaceType:
			builder.WriteString(fmt.Sprintf("%s%s %s.%s {\n", indentHead, InterfaceModifier, PkgNumberName, name))
			supers, methods := InterfaceFieldsStrings(typ, indentBody)
			if len(supers) > 0 {
				builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, SupperArea, strings.Join(supers, "")))
			}
			if len(methods) > 0 {
				builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, MethodArea, strings.Join(methods, "")))
			}
		default:
			builder.WriteString(fmt.Sprintf("%s%s %s.%s {\n", indentHead, ClassModifier, PkgNumberName, name))
			switch typ.(type) {
			case *ast.StructType:
				switch typ.(type) {
				case *ast.StructType:
					supers, fields := StructFieldsStrings(typ.(*ast.StructType), indentBody)
					if len(supers) > 0 {
						builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, SupperArea, strings.Join(supers, "")))
					}
					if len(fields) > 0 {
						builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, FieldArea, strings.Join(fields, "")))
					}
				}
				if len(dtyp.Methods) > 0 {
					builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, MethodArea, DocFuncsStrings(dtyp.Methods, indentBody, "")))
				}
				if len(dtyp.Funcs) > 0 {
					builder.WriteString(fmt.Sprintf("%s%s\n%s", indentBody, ConstructorArea, DocFuncsStrings(dtyp.Funcs, indentBody, "")))
				}
				if len(dtyp.Vars)+len(dtyp.Consts) > 0 {
					builder.WriteString(fmt.Sprintf("%s%s\n", indentBody, ScopeArea))
					if len(dtyp.Consts) > 0 {
						builder.WriteString(DocValuesStrings(dtyp.Consts, indentBody, "="))
					}
					if len(dtyp.Vars) > 0 {
						builder.WriteString(DocValuesStrings(dtyp.Vars, indentBody, ":"))
					}
				}
			default:
				builder.WriteString(fmt.Sprintf("%s%s\n", indentBody, PrototypeArea))
				builder.WriteString(fmt.Sprintf("%s%s\n", indentBody, types.ExprString(typ)))
			}
		}
		builder.WriteString(fmt.Sprintf("%s}\n", indentHead))
	}
	return builder.String()
}
