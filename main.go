package main

import (
	"github.com/VegetableDoggies/plantuml-gen-for-go/puml"
	"github.com/VegetableDoggies/plantuml-gen-for-go/utils/filepathx"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "rootdir",
		Aliases:  []string{"r", "R"},
		Usage:    "The path of project or package",
		Required: true,
	},
	&cli.StringSliceFlag{
		Name:    "excldirs",
		Aliases: []string{"e", "E"},
		Usage:   "The excluded dirs",
	},
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o", "O"},
		Usage:   "The output path",
	},
	&cli.BoolFlag{
		Name:    "isFlat",
		Aliases: []string{"f", "F"},
		Usage:   "default false, If true, make the packages flat",
		Value:   false,
	},
}

//go:generate go run main.go -r D:\Programs\Github组织\VegetableDoggies\plantuml-gen-for-go
func main() {
	app := &cli.App{
		Flags: flags,
		Name:  "go-puml-gen",
		Usage: "to generate .puml files",
		Action: func(cCtx *cli.Context) error {
			rootdir, err := filepath.Abs(cCtx.String("rootdir"))
			if err != nil {
				log.Fatal("\n[ERROR] Wrong rootdir :", rootdir, err)
			}
			excldirs := cCtx.StringSlice("excldirs")
			for i := range excldirs {
				excldir, err := filepath.Abs(excldirs[i])
				if err != nil {
					log.Fatal("\n[ERROR] Wrong excldir :", excldirs[i], err)
				}
				excldirs[i] = excldir
			}
			isFlat := cCtx.Bool("isFlat")
			output := cCtx.String("output")
			if output == "" || filepathx.IsDir(output) {
				if isFlat {
					output = filepath.Join(rootdir, filepath.Base(rootdir)+time.Now().Format("20060102150405")+"F.puml")
				} else {
					output = filepath.Join(rootdir, filepath.Base(rootdir)+time.Now().Format("20060102150405")+".puml")
				}
				log.Printf("The default output path is active : %s\n", output)
			}
			file, err := os.Create(output)
			if err != nil {
				log.Fatal("\n[ERROR] Can't Create the file :", output, err)
			}
			defer file.Close()
			if _, err = file.Write([]byte(puml.NewPortrait(rootdir, excldirs, isFlat).Puml)); err != nil {
				log.Fatal("\n[ERROR] Can't Write to file :", output, err)
			}
			log.Printf("Success!\n\n>>> Output = %s\n\n", output)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
