package models

type ValidFields struct {
	EntityName string
	Fields     map[string]FieldMetaData
}

func (f ValidFields) GetFieldType(s string) FieldType {
	fmd, ok := f.Fields[s]
	if !ok {
		return Undefined
	}
	return fmd.Type
}

func (f ValidFields) IsAnalyzed(s string) bool {
	fmd, ok := f.Fields[s]
	if !ok {
		return false
	}
	return fmd.IsAnalyzed
}
