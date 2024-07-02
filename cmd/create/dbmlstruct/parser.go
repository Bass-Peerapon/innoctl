package dbmlstruct

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Parser declaration
type Parser struct {
	s *Scanner

	// current token & literal
	token Token
	lit   string

	Debug bool
}

// NewParser ...
func NewParser(s *Scanner) *Parser {
	return &Parser{
		s:     s,
		token: ILLEGAL,
		lit:   "",
		Debug: os.Getenv("DBML_PARSER_DEBUG") == "true",
	}
}

// Parse ...
func (p *Parser) Parse() (*DBML, error) {
	dbml := &DBML{}
	for {
		p.next()
		switch p.token {
		case PROJECT:
			project, err := p.parseProject()
			if err != nil {
				return nil, err
			}
			p.debug("project", project)
			dbml.Project = *project
		case TABLE:
			table, err := p.parseTable()
			if err != nil {
				return nil, err
			}
			p.debug("table", table)

			// TODO:
			// * register table to tables map, for check ref
			dbml.Tables = append(dbml.Tables, *table)

		case REF:
			ref, err := p.parseRefs()
			if err != nil {
				return nil, err
			}
			p.debug("Refs", ref)

			// TODO:
			// * Check refs is valid or not (by tables map)
			dbml.Refs = append(dbml.Refs, *ref)

		case ENUM:
			enum, err := p.parseEnum()
			if err != nil {
				return nil, err
			}
			p.debug("Enum", enum)
			dbml.Enums = append(dbml.Enums, *enum)

		case TABLEGROUP:
			tableGroup, err := p.parseTableGroup()
			if err != nil {
				return nil, err
			}
			p.debug("TableGroup", tableGroup)
			dbml.TableGroups = append(dbml.TableGroups, *tableGroup)
		case EOF:
			return dbml, nil
		default:
			p.debug("token", p.token.String(), "lit", p.lit)
			return nil, p.expect("Project, Ref, Table, Enum, TableGroup")
		}
	}
}

func (p *Parser) parseTableGroup() (*TableGroup, error) {
	tableGroup := &TableGroup{}
	p.next()
	if p.token != IDENT && p.token != DSTRING {
		return nil, fmt.Errorf("TableGroup name is invalid: %s", p.lit)
	}
	tableGroup.Name = p.lit
	p.next()
	if p.token != LBRACE {
		return nil, p.expect("{")
	}
	p.next()

	for p.token == IDENT || p.token == DSTRING {
		tableGroup.Members = append(tableGroup.Members, p.lit)
		p.next()
	}
	if p.token != RBRACE {
		return nil, p.expect("}")
	}
	return tableGroup, nil
}

func (p *Parser) parseEnum() (*Enum, error) {
	enum := &Enum{}
	p.next()
	if !IsIdent(p.token) && p.token != DSTRING {
		return nil, fmt.Errorf("enum name is invalid: %s", p.lit)
	}
	enum.Name = p.lit
	p.next()
	if p.token != LBRACE {
		return nil, p.expect("{")
	}
	p.next()

	for IsIdent(p.token) {
		enumValue := EnumValue{
			Name: p.lit,
		}
		p.next()
		if p.token == LBRACK {
			// handle [Note: ...]
			p.next()
			if p.token == NOTE {
				note, err := p.parseDescription()
				if err != nil {
					return nil, p.expect("note: 'string'")
				}
				enumValue.Note = note
				p.next()
			}
			if p.token != RBRACK {
				return nil, p.expect("]")
			}
			p.next()
		}
		enum.Values = append(enum.Values, enumValue)
	}
	if p.token != RBRACE {
		return nil, p.expect("}")
	}
	return enum, nil
}

func (p *Parser) parseRefs() (*Ref, error) {
	ref := &Ref{}
	p.next()

	// Handle for Ref <optional_name>...
	if p.token == IDENT {
		ref.Name = p.lit
		p.next()
	}

	// Ref: from > to
	if p.token == COLON {
		p.next()
		rel, err := p.parseRelationship()
		if err != nil {
			return nil, err
		}
		ref.Relationships = append(ref.Relationships, *rel)
		return ref, nil
	}

	if p.token == LBRACE {
		p.next()

		for {
			if p.token == RBRACE {
				return ref, nil
			} else if p.token == IDENT || p.token == DSTRING {
				rel, err := p.parseRelationship()
				if err != nil {
					return nil, err
				}
				ref.Relationships = append(ref.Relationships, *rel)
			} else {
				return nil, p.expect("Ref: { from > to }")
			}
			p.next()
		}
	}

	return nil, p.expect("Ref: | Refs {}")
}

func (p *Parser) parseRelationship() (*Relationship, error) {
	rel := &Relationship{}
	if p.token != IDENT && p.token != DSTRING {
		return nil, p.expect("(rel from) table.column_name")
	}

	rel.From = p.lit

	p.next()
	if reltype, ok := RelationshipMap[p.token]; ok {
		rel.Type = reltype
	} else {
		return nil, p.expect("> | < | -")
	}

	p.next()
	if p.token != IDENT {
		return nil, p.expect("(rel to) table.column_name")
	}
	rel.To = p.lit
	return rel, nil
}

func (p *Parser) parseTable() (*Table, error) {
	table := &Table{}
	p.next()
	switch p.token {
	case IDENT, DSTRING:
		// pass
	default:
		if m, _ := regexp.MatchString("^[a-zA-Z1-9]+$", p.lit); !m {
			return nil, fmt.Errorf("table name is invalid: %s", p.lit)
		}
	}
	table.Name = p.lit

	p.next()

	switch p.token {
	case AS:
		// handle as
		p.next()
		switch p.token {
		case STRING, IDENT:
			table.As = p.lit
		default:
			return nil, p.expect("as NAME")
		}
		p.next()
		fallthrough
	case LBRACE:
		p.next()
		for {
			switch p.token {
			case INDEXES:
				indexes, err := p.parseIndexes()
				if err != nil {
					return nil, err
				}
				table.Indexes = indexes
			case RBRACE:
				return table, nil
			default:
				columnName := p.lit
				currentToken := p.token
				p.next()
				if currentToken == NOTE && p.token == COLON {
					note, err := p.parseString()
					if err != nil {
						return nil, err
					}
					table.Note = note
					p.next()
				} else {
					column, err := p.parseColumn(columnName)
					if err != nil {
						return nil, err
					}
					table.Columns = append(table.Columns, *column)
				}
			}
		}
	default:
		return nil, p.expect("{")
	}
}

func (p *Parser) parseIndexes() ([]Index, error) {
	indexes := []Index{}

	p.next()
	if p.token != LBRACE {
		return nil, p.expect("{")
	}

	p.next()
	for {
		if p.token == RBRACE {
			p.next() // pop }
			return indexes, nil
		}
		// parse an Index
		index, err := p.parseIndex()
		if err != nil {
			return nil, err
		}
		p.debug("index", index)
		indexes = append(indexes, *index)
	}
}

func (p *Parser) parseIndex() (*Index, error) {
	index := &Index{}

	if p.token == LPAREN {
		p.next()
		for IsIdent(p.token) {
			index.Fields = append(index.Fields, p.lit)
			p.next()
			if p.token == COMMA {
				p.next()
			}
		}
		if p.token != RPAREN {
			return nil, p.expect(")")
		}
	} else if IsIdent(p.token) {
		index.Fields = append(index.Fields, p.lit)
	} else {
		return nil, p.expect("field_name")
	}

	p.next()

	if p.token == LBRACK {
		// Handle index setting [settings...]
		commaAllowed := false

		for {
			p.next()
			switch {
			case p.token == IDENT && strings.ToLower(p.lit) == "name":
				name, err := p.parseDescription()
				if err != nil {
					return nil, p.expect("name: 'index_name'")
				}
				index.Settings.Name = name
			case p.token == NOTE:
				note, err := p.parseDescription()
				if err != nil {
					return nil, p.expect("note: 'index note'")
				}
				index.Settings.Note = note
			case p.token == PK:
				index.Settings.PK = true
			case p.token == UNIQUE:
				index.Settings.Unique = true
			case p.token == TYPE:
				p.next()
				if p.token != COLON {
					return nil, p.expect(":")
				}
				p.next()
				if p.token != IDENT || (p.lit != "hash" && p.lit != "btree") {
					return nil, p.expect("hash|btree")
				}
				index.Settings.Type = p.lit
			case p.token == COMMA:
				if !commaAllowed {
					return nil, p.expect("[index settings...]")
				}
			case p.token == RBRACK:
				p.next()
				return index, nil
			default:
				return nil, p.expect("note|name|type|pk|unique")
			}
			commaAllowed = !commaAllowed
		}
	}

	return index, nil
}

func (p *Parser) parseColumn(name string) (*Column, error) {
	column := &Column{
		Name: name,
	}
	if p.token != IDENT {
		return nil, p.expect("int, varchar,...")
	}
	column.Type = p.lit
	p.next()

	// parse for type
	switch p.token {
	case LPAREN:
		p.next()
		if p.token != INT {
			return nil, p.expect("int")
		}
		column.Type = fmt.Sprintf("%s(%s)", column.Type, p.lit)
		p.next()
		if p.token != RPAREN {
			return nil, p.expect(RPAREN.String())
		}
		p.next()
		if p.token != LBRACK {
			break
		}
		fallthrough
	case LBRACK:
		// handle parseColumn
		columnSetting, err := p.parseColumnSettings()
		if err != nil {
			return nil, fmt.Errorf("parse column settings: %w", err)
		}
		p.next() // remove ']'
		column.Settings = *columnSetting
	}

	p.debug("column", column)
	return column, nil
}

func (p *Parser) parseColumnSettings() (*ColumnSetting, error) {
	columnSetting := &ColumnSetting{Null: true}
	commaAllowed := false

	for {
		p.next()
		switch p.token {
		case PK:
			columnSetting.PK = true
		case PRIMARY:
			p.next()
			if p.token != KEY {
				return nil, p.expect("KEY")
			}
			columnSetting.PK = true
		case REF:
			p.next()
			if p.token != COLON {
				return nil, p.expect(":")
			}
			p.next()
			if p.token != LSS && p.token != GTR && p.token != SUB {
				return nil, p.expect("< | > | -")
			}
			columnSetting.Ref.Type = RelationshipMap[p.token]
			p.next()
			if p.token != IDENT {
				return nil, p.expect("table.column_id")
			}
			columnSetting.Ref.To = p.lit
		case NOT:
			p.next()
			if p.token != NULL {
				return nil, p.expect("null")
			}
			columnSetting.Null = false
		case UNIQUE:
			columnSetting.Unique = true
		case INCREMENT:
			columnSetting.Increment = true
		case DEFAULT:
			p.next()
			if p.token != COLON {
				return nil, p.expect(":")
			}
			p.next()
			switch p.token {
			case STRING, DSTRING, TSTRING, INT, FLOAT, EXPR:
				// TODO:
				//	* handle default value by expr
				//	* validate default value by type
				columnSetting.Default = p.lit
			default:
				return nil, p.expect("default value")
			}
		case NOTE:
			str, err := p.parseDescription()
			if err != nil {
				return nil, err
			}
			columnSetting.Note = str
		case COMMA:
			if !commaAllowed {
				return nil, p.expect("pk | primary key | unique")
			}
		case RBRACK:
			return columnSetting, nil
		default:
			return nil, p.expect("pk, primary key, unique")
		}
		commaAllowed = !commaAllowed
	}
}

func (p *Parser) parseProject() (*Project, error) {
	project := &Project{}
	p.next()
	if p.token != IDENT && p.token != DSTRING {
		return nil, p.expect("project_name")
	}

	project.Name = p.lit
	p.next()

	if p.token != LBRACE {
		return nil, p.expect("{")
	}
	for {
		p.next()
		switch p.token {
		case IDENT:
			switch p.lit {
			case "database_type":
				str, err := p.parseDescription()
				if err != nil {
					return nil, err
				}
				project.DatabaseType = str
			default:
				return nil, p.expect("database_type")
			}
		case NOTE:
			note, err := p.parseDescription()
			if err != nil {
				return nil, err
			}
			project.Note = note
		case RBRACE:
			return project, nil
		default:
			return nil, fmt.Errorf("invalid token %s", p.lit)
		}
	}
}

func (p *Parser) parseString() (string, error) {
	p.next()
	switch p.token {
	case STRING, DSTRING, TSTRING:
		return p.lit, nil
	default:
		return "", p.expect("string, double quote string, triple string")
	}
}

func (p *Parser) parseDescription() (string, error) {
	p.next()
	if p.token != COLON {
		return "", p.expect(":")
	}
	return p.parseString()
}

func (p *Parser) next() {
	for {
		p.token, p.lit = p.s.Read()
		// p.debug("token:", p.String(), "lit:", p.lit)
		if p.token != COMMENT {
			break
		}
	}
}

func (p *Parser) expect(expected string) error {
	l, c := p.s.LineInfo()
	return fmt.Errorf("[%d:%d] invalid token '%s', expected: '%s'", l, c, p.lit, expected)
}

func (p *Parser) debug(args ...interface{}) {
	if p.Debug {
		for _, arg := range args {
			fmt.Printf("%#v\t", arg)
		}
		fmt.Println()
	}
}
