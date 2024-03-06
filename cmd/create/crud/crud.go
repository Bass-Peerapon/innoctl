package crud

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"regexp"
	"strings"

	templater "github.com/Bass-Peerapon/innoctl/cmd/create/crud/template"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

var (
	outputFileName string
	fileName       string
)

type Model struct {
	Name           string
	CamelCase      string
	LowerCamelCase string
}

func NewModel(name string) *Model {
	return &Model{
		Name:           name,
		CamelCase:      strcase.ToCamel(name),
		LowerCamelCase: strcase.ToLowerCamel(name),
	}
}

// crudRepoCmd represents the crudRepo command
var CrudRepoCmd = &cobra.Command{
	Use:   "crud",
	Short: "generate crud in repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := cmd.Flags().GetString("db_driver")
		if err != nil {
			fmt.Println(err)
			return
		}
		switch db {
		case "postgres":
			all, err := cmd.Flags().GetBool("all")
			if err != nil {
				fmt.Println(err)
				return
			}
			if all {
				if err := GenReadPostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}
				if err := GenCreatePostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}
				if err := GenUpdatePostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}
				if err := GenDeletePostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}

				return
			}
			flagRead, err := cmd.Flags().GetBool("read")
			if err != nil {
				fmt.Println(err)
				return
			}
			if flagRead {
				if err := GenReadPostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}
			}
			flagCreate, err := cmd.Flags().GetBool("create")
			if err != nil {
				fmt.Println(err)
				return
			}
			if flagCreate {
				if err := GenCreatePostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}
			}
			flagUpdate, err := cmd.Flags().GetBool("update")
			if err != nil {
				fmt.Println(err)
				return
			}
			if flagUpdate {
				if err := GenUpdatePostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}
			}
			flagDelete, err := cmd.Flags().GetBool("delete")
			if err != nil {
				fmt.Println(err)
				return
			}
			if flagDelete {
				if err := GenDeletePostgres(fileName, outputFileName); err != nil {
					fmt.Println(err)
				}
			}
		case "mongodb":

		default:
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crudRepoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	CrudRepoCmd.Flags().StringVarP(&fileName, "file", "f", "", "path file model")
	if err := CrudRepoCmd.MarkFlagRequired("file"); err != nil {
		fmt.Println(err)
	}
	CrudRepoCmd.Flags().StringVarP(&outputFileName, "output", "o", "", "path file repository")
	if err := CrudRepoCmd.MarkFlagRequired("output"); err != nil {
		fmt.Println(err)
	}
	CrudRepoCmd.Flags().StringP("db_driver", "", "postgres", "select db_driver (postgres , mongodb) defult postgres")
	CrudRepoCmd.Flags().BoolP("all", "a", false, "generate sql script (insert ,query ,update ,delete)")
	CrudRepoCmd.Flags().BoolP("create", "c", false, "generate sql script insert")
	CrudRepoCmd.Flags().BoolP("read", "r", false, "generate sql script query")
	CrudRepoCmd.Flags().BoolP("update", "u", false, "generate sql script update")
	CrudRepoCmd.Flags().BoolP("delete", "d", false, "generate sql script delete")
}

func getTag(s string) (string, error) {
	tag := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(s), "db:", ""), `"`, "")
	if tag == "" {
		return "", errors.New("not found tag `db`")
	}
	return tag, nil
}

func GenCreatePostgres(fn string, fnOut string) error {
	m, tagsDB, err := readModel(fn)
	if err != nil {
		return err
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)

	// tmp := strings.Join([]string{string(fo), string(templater.PostgresCreate)}, "\n")
	tmp := bytes.Join([][]byte{fo, templater.PostgresCreate}, []byte("\n"))

	addFunc := func(x, y int) int {
		return x + y
	}
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
		"add":       addFunc,
	}
	t, err := template.New("create").Funcs(funcMap).Parse(string(tmp))
	if err != nil {
		return err
	}
	if err = t.Execute(&buf, e); err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}
	return nil
}

func GenReadPostgres(fn string, fnOut string) error {
	m, tagsDB, err := readModel(fn)
	if err != nil {
		return err
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)
	tmp := bytes.Join([][]byte{fo, templater.PosrgresRead}, []byte("\n"))
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
	}
	t, err := template.New("read").Funcs(funcMap).Parse(string(tmp))
	if err != nil {
		return err
	}
	if err = t.Execute(&buf, e); err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}

	return nil
}

func GenUpdatePostgres(fn string, fnOut string) error {
	m, tagsDB, err := readModel(fn)
	if err != nil {
		return err
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)
	tmp := bytes.Join([][]byte{fo, templater.PostgresUpdate}, []byte("\n"))

	addFunc := func(x, y int) int {
		return x + y
	}
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
		"add":       addFunc,
	}
	t, err := template.New("update").Funcs(funcMap).Parse(string(tmp))
	if err != nil {
		return err
	}
	if err = t.Execute(&buf, e); err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}

	return nil
}

func GenDeletePostgres(fn string, fnOut string) error {
	m, tagsDB, err := readModel(fn)
	if err != nil {
		return err
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)
	tmp := bytes.Join([][]byte{fo, templater.PostgresDelete}, []byte("\n"))
	addFunc := func(x, y int) int {
		return x + y
	}
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
		"add":       addFunc,
	}
	t, err := template.New("update").Funcs(funcMap).Parse(string(tmp))
	if err != nil {
		return err
	}
	if err = t.Execute(&buf, e); err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}

	return nil
}

func readModel(fn string) (*Model, []string, error) {
	var m *Model
	var tagsDB []string

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fn, nil, 0)
	if err != nil {
		return nil, nil, err
	}
	reg := regexp.MustCompile(`db:"(.+)" `)
	for range f.Scope.Objects {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				if decl.Tok == token.IMPORT {
					continue
				}
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						m = NewModel(spec.Name.String())
						switch st := spec.Type.(type) {
						case *ast.StructType:
							for _, field := range st.Fields.List {
								tag := field.Tag.Value
								tags := reg.FindAllString(tag, 1)
								switch field.Type.(type) {
								case *ast.StarExpr:
									t, err := getTag(tags[0])
									if err != nil {
										fmt.Println(err.Error())
										continue
									}
									if t != "-" && t != "" {
										tagsDB = append(tagsDB, t)
									}
								case *ast.ArrayType:
									arr := field.Type.(*ast.ArrayType)
									if arr.Lbrack.IsValid() {
										switch arr.Elt.(type) {
										case *ast.StarExpr:
											t, err := getTag(tags[0])
											if err != nil {
												fmt.Println(err.Error())
												continue
											}
											if t != "-" && t != "" {
												tagsDB = append(tagsDB, t)
											}
										case *ast.Ident:
											t, err := getTag(tags[0])
											if err != nil {
												fmt.Println(err.Error())
												continue
											}
											if t != "-" && t != "" {
												tagsDB = append(tagsDB, t)
											}
										}
									}

								case *ast.Ident:
									t, err := getTag(tags[0])
									if err != nil {
										fmt.Println(err.Error())
										continue
									}
									if t != "-" && t != "" {
										tagsDB = append(tagsDB, t)
									}
								case *ast.SelectorExpr:
									t, err := getTag(tags[0])
									if err != nil {
										fmt.Println(err.Error())
										continue
									}
									if t != "-" && t != "" {
										tagsDB = append(tagsDB, t)
									}

								}
							}
						}
					}
				}
			}
		}
		break
	}
	return m, tagsDB, nil
}

func getModelFromRepo(fn string) (*Model, error) {
	var m *Model

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fn, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.IMPORT {
				continue
			}
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					m = NewModel(spec.Name.String())
					return m, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("not found struct")
}
