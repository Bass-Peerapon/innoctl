package dbmlstruct

import (
	"fmt"
	"regexp"
	"strings"
)

type Opts struct {
	From          string
	Out           string
	Package       string
	FieldTags     []string
	RememberAlias bool
	Recursive     bool
	Exclude       string
}

// Generate go model
func Generate(opts Opts) {
	var pattern *regexp.Regexp
	if strings.TrimSpace(opts.Exclude) != "" {
		pattern, _ = regexp.Compile(opts.Exclude)
	}

	dbmls := parseDBML(opts.From, opts.Recursive, pattern)

	g := newgen()
	g.out = opts.Out
	g.gopackage = opts.Package
	g.fieldtags = opts.FieldTags

	for _, dbml := range dbmls {
		g.reset(opts.RememberAlias)
		g.dbml = dbml
		if err := g.generate(); err != nil {
			fmt.Printf("Error generate file %s", err)
		}
	}
}

func GenerateFormString(r string) (string, error) {
	dbmls, err := parseDBMLFromString(r)
	if err != nil {
		return "", err
	}
	g := newgen()
	g.fieldtags = []string{"json", "db", "type"}
	strs := []string{}
	for _, dbml := range dbmls {
		g.dbml = dbml
		str, err := g.printGenerate()
		if err != nil {
			return "", err
		}
		strs = append(strs, str)
	}
	return strings.Join(strs, "\n"), nil
}
