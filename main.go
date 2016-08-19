package main

import (
	"bytes"
	"flag"
	"go/build"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	tmpl = `// This file is generated and should not be modified directly.
// Regenerate using ps2bs: 'go get github.com/codemodus/ps2bs; ps2bs'

{{ if .IsLibrary }}// Package {{ .PkgName }} provides functions which return pointers to
// built-ins and other commonly used standard library types.{{ end }}
package {{ .PkgName }}

{{- range $val := .Types }}
{{ if $.IsExported }}
// To{{ $val | TypeName | TitleCase }} returns a pointer to the type '{{ $val }}'.
{{- end }}
func {{ if $.IsExported }}To{{ else }}to{{ end }}{{ $val | TypeName | TitleCase }}({{ $val | TypeName | FirstAsLC }} {{ $val }}) *{{ $val }} {
	return &{{ $val | TypeName | FirstAsLC }}
}{{ end }}
`
)

var (
	types = []string{
		"bool",
		"byte",
		"complex128",
		"complex64",
		"error",
		"float32",
		"float64",
		"int",
		"int16",
		"int32",
		"int64",
		"int8",
		"rune",
		"string",
		"uint",
		"uint16",
		"uint32",
		"uint64",
		"uint8",
		"uintptr",
		"time.Time",
	}

	fileName = "ps2bs_gen.go"
)

type options struct {
	directory string
	exported  bool
}

type tmplContext struct {
	PkgName    string
	IsLibrary  bool
	IsExported bool
	Types      []string
}

func newTmplContext(opts *options) (*tmplContext, error) {
	destPkgName, isExisting, err := pkgInfo(opts.directory)
	if err != nil {
		return nil, err
	}

	ctx := &tmplContext{
		PkgName:    destPkgName,
		IsLibrary:  !isExisting,
		IsExported: opts.exported,
		Types:      types,
	}

	return ctx, nil
}

func execTmpl(ctx *tmplContext) ([]byte, error) {
	fs := funcMap()

	bb := &bytes.Buffer{}
	t, err := template.New("").Funcs(fs).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	if err = t.Execute(bb, ctx); err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

func main() {
	log.SetFlags(0)

	o := &options{}
	{
		flag.StringVar(&o.directory, "dir", ".", "destination directory")
		flag.BoolVar(&o.exported, "e", false, "export functions")
	}
	flag.Parse()

	if err := run(o); err != nil {
		log.Fatalln(err)
	}
}

func run(o *options) error {
	ctx, err := newTmplContext(o)
	if err != nil {
		return err
	}

	tbs, err := execTmpl(ctx)
	if err != nil {
		return err
	}

	bs, err := format.Source(tbs)
	if err != nil {
		return err
	}

	var aErr error
	func() {
		var f *os.File
		f, aErr = os.Create(filepath.Join(o.directory, fileName))
		if aErr != nil {
			return
		}

		defer func() {
			if aErr = f.Close(); aErr != nil {
				return
			}
		}()

		if _, aErr = f.Write(bs); aErr != nil {
			return
		}
	}()

	return aErr
}

func pkgInfo(dir string) (name string, isExisting bool, err error) {
	fs, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return "", false, err
	}

	if len(fs) == 0 || len(fs) == 1 && filepath.Base(fs[0]) == fileName {
		return filepath.Base(dir), false, nil
	}

	pkg, err := build.ImportDir(dir, 0)
	if err != nil {
		return "", false, err
	}

	return pkg.Name, true, nil
}

func firstAsLC(s string) string {
	return strings.ToLower(s[0:1])
}

func typeName(s string) string {
	ss := strings.Split(s, ".")
	return ss[len(ss)-1]
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"TitleCase": strings.Title,
		"FirstAsLC": firstAsLC,
		"TypeName":  typeName,
	}
}
