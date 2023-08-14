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
		Aliases:  []string{"d"},
		Usage:    "The path of project or package",
		Required: true,
	},
	&cli.StringSliceFlag{
		Name:    "excldirs",
		Aliases: []string{"e"},
		Usage:   "The excluded dirs",
	},
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "The output path",
	},
}

//go:generate go run xxx -d D:\Programs\go-puml-gen -e D:\Programs\go-puml-gen\utils
func main() {
	app := &cli.App{
		Flags: flags,
		Name:  "go-puml-gen",
		Usage: "to generate .puml files",
		Action: func(cCtx *cli.Context) error {
			rootdir, err := filepath.Abs(cCtx.String("rootdir"))
			if err != nil {
				log.Fatal("[ERROR] Wrong rootdir :", rootdir, err)
			}
			excldirs := cCtx.StringSlice("excldirs")
			for i := range excldirs {
				excldir, err := filepath.Abs(excldirs[i])
				if err != nil {
					log.Fatal("[ERROR] Wrong excldir :", excldirs[i], err)
				}
				excldirs[i] = excldir
			}
			p := puml.NewPortrait(rootdir, excldirs)
			output := cCtx.String("output")
			if output == "" || filepathx.IsDir(output) {
				output = filepath.Join(rootdir, filepath.Base(rootdir)+time.Now().Format("20060102150405")+".puml")
				log.Printf("The default output path is active : %s\n", output)
			}
			file, err := os.Create(output)
			if err != nil {
				log.Fatal("[ERROR] Can't Create the file :", output, err)
			}
			defer file.Close()
			if _, err = file.Write([]byte(p.Puml)); err != nil {
				log.Fatal("[ERROR] Can't Write to file :", output, err)
			}
			log.Printf("Success!\n\n>>> Output = %s\n\n", output)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
