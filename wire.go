package dynago

import "strconv"

func wireEncode(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return &wireString{v}
	case int:
		return &wireNumber{strconv.Itoa(v)}
	case int64:
		return &wireNumber{strconv.FormatInt(v, 10)}
	case bool:
		return &wireBool{v}
	case float64:
		return &wireNumber{strconv.FormatFloat(v, 'g', -1, 64)}
	case Number:
		return &wireNumber{string(v)}
	case Document:
		output := make(map[string]interface{}, len(v))
		for key, val := range v {
			output[key] = wireEncode(val)
		}
		return &wireMap{output}
	case []byte:
		return &wireBinary{v}
	case []string:
		return &wireStringSet{v}
	case StringSet:
		return &wireStringSet{v}
	case NumberSet:
		return &wireNumberSet{v}
	case BinarySet:
		return &wireBinarySet{v}
	case List:
		encList := make([]interface{}, len(v))
		for i, raw := range v {
			encList[i] = wireEncode(raw)
		}
		return &wireList{encList}
	case map[string]interface{}:
		return wireEncode(Document(v))
	default:
		panic(v)
	}
}

type wireString struct {
	S string
}

type wireStringSet struct {
	SS []string
}

type wireNumber struct {
	N string
}

type wireNumberSet struct {
	NS []string
}

type wireBool struct {
	BOOL bool `json:"BOOL"`
}

// Bonus! we don't have to do anything at all for binary. base64 is done by encoding/json.
type wireBinary struct {
	B []byte
}

type wireBinarySet struct {
	BS [][]byte
}

type wireList struct {
	L []interface{}
}

type wireMap struct {
	M map[string]interface{}
}

func wireDecode(original interface{}) interface{} {
	vv, ok := original.(map[string]interface{})
	if !ok {
		panic(original) // XXX DEBUG TODO
	}
	for typeCode, val := range vv {
		// TODO decode all the various types coming back.
		switch typeCode {
		case "S", "BOOL":
			return val
		case "NULL":
			return nil
		case "N":
			return Number(val.(string))
		case "NS":
			return wireDecodeNumberSet(val)
		case "SS":
			return wireDecodeStringSet(val)
		case "L":
			return wireDecodeList(val)
		}
	}
	return nil
}

func wireDecodeNumberSet(val interface{}) interface{} {
	return Number(val.(string))
}

func wireDecodeStringSet(val interface{}) interface{} {
	valSlice := val.([]interface{})
	resultSlice := StringSet{}
	for _, v := range valSlice {
		resultSlice = append(resultSlice, v.(string))
	}
	return resultSlice
}

func wireDecodeList(val interface{}) interface{} {
	return nil // TODO
}
