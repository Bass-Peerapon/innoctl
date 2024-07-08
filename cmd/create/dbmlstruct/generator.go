package dbmlstruct

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

type generator struct {
	dbml             *DBML
	out              string
	gopackage        string
	fieldtags        []string
	types            map[string]jen.Code
	shouldGenTblName bool
}

func newgen() *generator {
	return &generator{
		types: make(map[string]jen.Code),
	}
}

func (g *generator) reset(rememberAlias bool) {
	g.dbml = nil
	if !rememberAlias {
		g.types = make(map[string]jen.Code)
	}
}

func (g *generator) file() *jen.File {
	return jen.NewFilePathName(g.out, g.gopackage)
}

func (g *generator) generate() error {
	if err := g.genEnums(); err != nil {
		return err
	}
	return nil
}

func (g *generator) printGenerate() (string, error) {
	return g.printGenEnums()
}

func (g *generator) genEnums() error {
	for _, enum := range g.dbml.Enums {
		if err := g.genEnum(enum); err != nil {
			return err
		}
	}
	for _, table := range g.dbml.Tables {
		if err := g.genTable(table); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) printGenEnums() (string, error) {
	strs := []string{}
	for _, enum := range g.dbml.Enums {
		str, err := g.printGenEnum(enum)
		if err != nil {
			return "", err
		}
		strs = append(strs, str)

	}
	for _, table := range g.dbml.Tables {
		str, err := g.printGenTable(table)
		if err != nil {
			return "", err
		}
		strs = append(strs, str)
	}

	return strings.Join(strs, "\n"), nil
}

func (g *generator) genEnum(enum Enum) error {
	f := jen.NewFilePathName(g.out, g.gopackage)

	enumGoTypeName := NormalizeGoTypeName(enum.Name)

	f.Type().Id(enumGoTypeName).Int()

	f.Const().DefsFunc(func(group *jen.Group) {
		group.Id("_").Id(enumGoTypeName).Op("=").Iota()
		for _, value := range enum.Values {
			v := group.Id(NormalLizeGoName(value.Name))
			if value.Note != "" {
				v.Comment(value.Note)
			}
		}
	})

	g.types[enum.Name] = jen.Id(enumGoTypeName)

	return f.Save(fmt.Sprintf("%s/%s.enum.go", g.out, Normalize(enum.Name)))
}

func (g *generator) printGenEnum(enum Enum) (string, error) {
	enumGoTypeName := NormalizeGoTypeName(enum.Name)

	f := jen.Type().Id(enumGoTypeName).Int()

	f.Const().DefsFunc(func(group *jen.Group) {
		group.Id("_").Id(enumGoTypeName).Op("=").Iota()
		for _, value := range enum.Values {
			v := group.Id(NormalLizeGoName(value.Name))
			if value.Note != "" {
				v.Comment(value.Note)
			}
		}
	})

	g.types[enum.Name] = jen.Id(enumGoTypeName)

	return f.GoString(), nil
}

func (g *generator) genTable(table Table) error {
	f := jen.NewFilePathName(g.out, g.gopackage)

	tableGoTypeName := NormalizeGoTypeName(table.Name)

	f.Type().Id(tableGoTypeName).StructFunc(func(group *jen.Group) {
		tb := group.Id("TableName").Add(jen.Struct())
		dbTags := map[string]string{
			"json": "-",
			"db":   table.Name,
		}
		for _, column := range table.Columns {
			columnName := NormalLizeGoName(column.Name)
			if column.Settings.PK {
				dbTags["pk"] = columnName
			}
			columnOriginName := Normalize(column.Name)
			typ, ok := g.getJenType(column.Type)
			if !ok {
				typ = jen.Qual("constants", column.Type)
			}
			if column.Settings.Note != "" {
				group.Comment(column.Settings.Note)
			}

			gotags := make(map[string]string)
			for _, t := range g.fieldtags {
				if t != "type" {
					gotags[strings.TrimSpace(t)] = columnOriginName
					continue
				}
				gtype := g.getType(column.Type)
				if gtype == "" {
					gtype = column.Type
				}
				gotags[strings.TrimSpace(t)] = gtype
			}
			group.Id(columnName).Add(typ).Tag(gotags)

		}
		tb.Tag(dbTags)
	})

	return f.Save(fmt.Sprintf("%s/%s.table.go", g.out, Normalize(table.Name)))
}

func (g *generator) printGenTable(table Table) (string, error) {
	tableGoTypeName := NormalizeGoTypeName(table.Name)

	f := jen.Type().Id(tableGoTypeName).StructFunc(func(group *jen.Group) {
		tb := group.Id("TableName").Add(jen.Struct())
		dbTags := map[string]string{
			"json": "-",
			"db":   table.Name,
		}
		for _, column := range table.Columns {
			columnName := NormalLizeGoName(column.Name)
			if column.Settings.PK {
				dbTags["pk"] = columnName
			}

			columnOriginName := Normalize(column.Name)
			typ, ok := g.getJenType(column.Type)
			if !ok {
				typ = jen.Qual("constants", column.Type)
			}
			if column.Settings.Note != "" {
				group.Comment(column.Settings.Note)
			}

			gotags := make(map[string]string)
			for _, t := range g.fieldtags {
				if t != "type" {
					gotags[strings.TrimSpace(t)] = columnOriginName
					continue
				}
				gtype := g.getType(column.Type)
				if gtype == "" {
					gtype = column.Type
				}
				gotags[strings.TrimSpace(t)] = gtype
			}
			group.Id(columnName).Add(typ).Tag(gotags)
		}

		tb.Tag(dbTags)
	})

	return f.GoString(), nil
}

const primeTypePattern = `^(\w+)(\(d+\))?`

var (
	regexType    = regexp.MustCompile(primeTypePattern)
	builtinTypes = map[string]jen.Code{
		"int":       jen.Int(),
		"int8":      jen.Int8(),
		"int16":     jen.Int16(),
		"int32":     jen.Int32(),
		"int64":     jen.Int64(),
		"integer":   jen.Int(),
		"bigint":    jen.Int64(),
		"uint":      jen.Uint(),
		"uint8":     jen.Uint8(),
		"uint16":    jen.Uint16(),
		"uint32":    jen.Uint32(),
		"uint64":    jen.Uint64(),
		"float":     jen.Float64(),
		"float32":   jen.Float32(),
		"float64":   jen.Float64(),
		"numeric":   jen.Float64(),
		"bool":      jen.Bool(),
		"text":      jen.String(),
		"varchar":   jen.String(),
		"char":      jen.String(),
		"byte":      jen.Byte(),
		"rune":      jen.Rune(),
		"timestamp": jen.Op("*").Qual("time", "Time"),
		"time":      jen.Op("*").Qual("time", "Time"),
		"datetime":  jen.Op("*").Qual("time", "Time"),
		"date":      jen.Op("*").Qual("time", "Time"),
		"uuid":      jen.Op("*").Qual("github.com/gofrs/uuid", "UUID"),
	}
	builtinTypesGo = map[string]string{
		"int":       "int",
		"int8":      "int8",
		"int16":     "int16",
		"int32":     "int32",
		"int64":     "int64",
		"integer":   "int",
		"bigint":    "int64",
		"uint":      "uint",
		"uint8":     "uint8",
		"uint16":    "uint16",
		"uint32":    "uint32",
		"uint64":    "uint64",
		"float":     "float64",
		"float32":   "float32",
		"float64":   "float64",
		"numeric":   "float64",
		"bool":      "bool",
		"text":      "string",
		"varchar":   "string",
		"char":      "string",
		"byte":      "byte",
		"rune":      "rune",
		"uuid":      "uuid",
		"timestamp": "timestamp",
		"time":      "timestamp",
		"datetime":  "date",
		"date":      "date",
	}
)

func (g *generator) getJenType(s string) (jen.Code, bool) {
	m := regexType.FindStringSubmatch(s)
	if len(m) >= 2 {
		// lookup for builtin type
		if t, ok := builtinTypes[strings.ToLower(m[1])]; ok {
			return t, ok
		}
	}
	typ, ok := g.types[s]
	return typ, ok
}

func (g *generator) getType(s string) string {
	m := regexType.FindStringSubmatch(s)
	if len(m) >= 2 {
		// lookup for builtin type
		if t, ok := builtinTypesGo[strings.ToLower(m[1])]; ok {
			return t
		}
	}
	return ""
}
