package dbmlstruct

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOF-1]
	_ = x[COMMENT-2]
	_ = x[_literalBeg-3]
	_ = x[IDENT-4]
	_ = x[INT-5]
	_ = x[FLOAT-6]
	_ = x[IMAG-7]
	_ = x[STRING-8]
	_ = x[DSTRING-9]
	_ = x[TSTRING-10]
	_ = x[EXPR-11]
	_ = x[_literalEnd-12]
	_ = x[_operatorBeg-13]
	_ = x[SUB-14]
	_ = x[LSS-15]
	_ = x[GTR-16]
	_ = x[LPAREN-17]
	_ = x[LBRACK-18]
	_ = x[LBRACE-19]
	_ = x[COMMA-20]
	_ = x[PERIOD-21]
	_ = x[RPAREN-22]
	_ = x[RBRACK-23]
	_ = x[RBRACE-24]
	_ = x[SEMICOLON-25]
	_ = x[COLON-26]
	_ = x[_operatorEnd-27]
	_ = x[_keywordBeg-28]
	_ = x[PROJECT-29]
	_ = x[TABLE-30]
	_ = x[ENUM-31]
	_ = x[REF-32]
	_ = x[AS-33]
	_ = x[TABLEGROUP-34]
	_ = x[_keywordEnd-35]
	_ = x[_miscBeg-36]
	_ = x[PRIMARY-37]
	_ = x[KEY-38]
	_ = x[PK-39]
	_ = x[NOTE-40]
	_ = x[UNIQUE-41]
	_ = x[NOT-42]
	_ = x[NULL-43]
	_ = x[INCREMENT-44]
	_ = x[DEFAULT-45]
	_ = x[INDEXES-46]
	_ = x[TYPE-47]
	_ = x[DELETE-48]
	_ = x[UPDATE-49]
	_ = x[NO-50]
	_ = x[ACTION-51]
	_ = x[RESTRICT-52]
	_ = x[SET-53]
	_ = x[_miscEnd-54]
}

const _Token_name = "ILLEGALEOFCOMMENT_literalBegIDENTINTFLOATIMAGSTRINGDSTRINGTSTRINGEXPR_literalEnd_operatorBegSUBLSSGTRLPARENLBRACKLBRACECOMMAPERIODRPARENRBRACKRBRACESEMICOLONCOLON_operatorEnd_keywordBegPROJECTTABLEENUMREFASTABLEGROUP_keywordEnd_miscBegPRIMARYKEYPKNOTEUNIQUENOTNULLINCREMENTDEFAULTINDEXESTYPEDELETEUPDATENOACTIONRESTRICTSET_miscEnd"

var _Token_index = [...]uint16{0, 7, 10, 17, 28, 33, 36, 41, 45, 51, 58, 65, 69, 80, 92, 95, 98, 101, 107, 113, 119, 124, 130, 136, 142, 148, 157, 162, 174, 185, 192, 197, 201, 204, 206, 216, 227, 235, 242, 245, 247, 251, 257, 260, 264, 273, 280, 287, 291, 297, 303, 305, 311, 319, 322, 330}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
