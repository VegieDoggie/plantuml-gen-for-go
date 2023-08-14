# plantuml-gen-for-go

这是一个将go项目或者go模块转换为PlanUML类图的工具，以图形的方式概览源码。

## 背景

2021年，我为了研读`go-etheruem`源码，于是编写了自动将源码生成UML类图的工具，这样就能大大加快源码阅读速度! 

## 预览
> (默认)嵌套式的UML类图 (说明: 由于渲染的原因，同一个包会显示两层)

![默认格式](https://github.com/VegetableDoggies/plantuml-gen-for-go/images/plantuml-gen-go20230814190411-0.png)

> (扁平化)所有包同一层级展示的UML类图

![扁平化格式](https://github.com/VegetableDoggies/plantuml-gen-for-go/images/plantuml-gen-go20230814190429F-0.png)

## 快速入门
> 说明: 想查看类图，需要IDE安装PlantUML Integration插件，地址: https://plugins.jetbrains.com/plugin/7017-plantuml-integration


```go
// 1- 下载并安装plantuml-gen-for-go
go install github.com/VegetableDoggies/plantuml-gen-for-go@1.1.0
// 2- 将指定目录的go项目或模块转换为UML类图
plantuml-gen-for-go -r D:/Programs/Github组织/VegetableDoggies/plantuml-gen-go
// 或 plantuml-gen-for-go -r D:/Programs/Github组织/VegetableDoggies/plantuml-gen-go -f
```

```cmd
NAME:
   go-puml-gen - to generate .puml files

USAGE:
   go-puml-gen [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rootdir value, -r value, -R value                                            The path of project or package
   --excldirs value, -e value, -E value [ --excldirs value, -e value, -E value ]  The excluded dirs
   --output value, -o value, -O value                                             The output path
   --isFlat, -f, -F                                                               default false, If true, make the packages flat (default: false)
   --help, -h                                                                     show help
```
