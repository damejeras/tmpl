package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

var stderr = log.New(os.Stderr, "", 0)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		stderr.Println("Can not read stdin.")
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		stderr.Println("The command is intended to work with pipes.")
		stderr.Println("Usage: tmpl values.json < template.tmpl")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		stderr.Println("Values file is required.")
		stderr.Println("Usage: tmpl values.json < template.tmpl")
		os.Exit(1)
	}

	valuesFile, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0644)
	if err != nil {
		stderr.Printf("Can not open file %q.\n", os.Args[1])
		panic(err)
	}

	defer valuesFile.Close()

	values := make(map[string]interface{})
	err = json.NewDecoder(valuesFile).Decode(&values)
	if err != nil {
		stderr.Printf("Can not decode %q.\n", os.Args[1])
		panic(err)
	}

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		stderr.Println("Can not read stdin.")
		panic(err)
	}

	tmpl, err := template.New("default").Parse(string(stdin))
	if err != nil {
		stderr.Println("Can not parse template from stdin.")
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, values); err != nil {
		stderr.Println("Can not execute template.")
		panic(err)
	}
}
