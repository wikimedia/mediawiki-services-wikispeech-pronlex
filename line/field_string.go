// Code generated by "stringer -type=Field"; DO NOT EDIT

package line

import "fmt"

const _Field_name = "OrthPosMorphDecompWordLangTrans1Translang1Trans2Translang2Trans3Translang3Trans4Translang4Trans5Translang5Trans6Translang6LemmaInflectionRule"

var _Field_index = [...]uint8{0, 4, 7, 12, 18, 26, 32, 42, 48, 58, 64, 74, 80, 90, 96, 106, 112, 122, 127, 141}

func (i Field) String() string {
	if i < 0 || i >= Field(len(_Field_index)-1) {
		return fmt.Sprintf("Field(%d)", i)
	}
	return _Field_name[_Field_index[i]:_Field_index[i+1]]
}
