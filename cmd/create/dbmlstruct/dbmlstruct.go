/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package dbmlstruct

import (
	"github.com/spf13/cobra"
)

var (
	from             = "database.dbml"
	out              = "model"
	gopackage        = "model"
	fieldtags        = []string{"json", "db", "type"}
	shouldGenTblName = false
	rememberAlias    = false
	recursive        = false
	exclude          = ""
)

var DbmlstructCmd = &cobra.Command{
	Use:   "dbml2struct",
	Short: "generate struct from dbml",
	Long:  `generate struct from dbml`,
	Run: func(cmd *cobra.Command, args []string) {
		Generate(Opts{
			From:          from,
			Out:           out,
			Package:       gopackage,
			FieldTags:     fieldtags,
			RememberAlias: rememberAlias,
			Recursive:     recursive,
			Exclude:       exclude,
		})
	},
}

func init() {
	flags := DbmlstructCmd.Flags()
	flags.StringVarP(&from, "from", "f", from, "source of dbml, can be https://dbdiagram.io/... | fire_name.dbml")
	flags.StringVarP(&out, "out", "o", out, "output folder")
	flags.StringVarP(&gopackage, "package", "p", gopackage, "single for multiple files")
	flags.StringArrayVarP(&fieldtags, "fieldtags", "t", fieldtags, "go field tags")
	flags.BoolVarP(&shouldGenTblName, "gen-table-name", "", shouldGenTblName, "should generate \"TableName\" function")
	flags.BoolVarP(&rememberAlias, "remember-alias", "", rememberAlias, "should remember table alias. Only applied if \"from\" is a directory")
	flags.BoolVarP(&recursive, "recursive", "", recursive, "recursive search directory. Only applied if \"from\" is a directory")
	flags.StringVarP(&exclude, "exclude", "E", exclude, "regex for exclude \"from\" files. Only applied if \"from\" is a directory")
}
