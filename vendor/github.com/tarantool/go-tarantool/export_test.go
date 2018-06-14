package tarantool

func (schema *Schema) ResolveSpaceIndex(s interface{}, i interface{}) (spaceNo, indexNo uint32, err error) {
	return schema.resolveSpaceIndex(s, i)
}
