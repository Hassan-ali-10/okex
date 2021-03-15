package helpers


//AlterInterfaceTypes "Alters interface types for int types converted to float and wraps it up into a map again with same keys"
func AlterInterfaceTypesToFloat(mapp map[string]interface{}) map[string]interface{} {

	var alteredMapp = make(map[string]interface{})

	for key, value := range mapp {

		switch v := value.(type) {
		case nil:
		case int:
			alteredMapp[key] = float64(v)
		case int16:
			alteredMapp[key] = float64(v)
		case int32:
			alteredMapp[key] = float64(v)
		case int64:
			alteredMapp[key] = float64(v)
		default:
			alteredMapp[key] = v
		}
	}

	return alteredMapp

}