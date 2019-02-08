// Code generated by "enumer -type=FieldOrder -json=true -transform=comment"; DO NOT EDIT.

package vtil

import (
	"encoding/json"
	"fmt"
)

const _FieldOrderName = "unknownprogressivettbbtbbt"

var _FieldOrderIndex = [...]uint8{0, 7, 18, 20, 22, 24, 26}

func (i FieldOrder) String() string {
	if i < 0 || i >= FieldOrder(len(_FieldOrderIndex)-1) {
		return fmt.Sprintf("FieldOrder(%d)", i)
	}
	return _FieldOrderName[_FieldOrderIndex[i]:_FieldOrderIndex[i+1]]
}

var _FieldOrderValues = []FieldOrder{0, 1, 2, 3, 4, 5}

var _FieldOrderNameToValueMap = map[string]FieldOrder{
	_FieldOrderName[0:7]:   0,
	_FieldOrderName[7:18]:  1,
	_FieldOrderName[18:20]: 2,
	_FieldOrderName[20:22]: 3,
	_FieldOrderName[22:24]: 4,
	_FieldOrderName[24:26]: 5,
}

// FieldOrderString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func FieldOrderString(s string) (FieldOrder, error) {
	if val, ok := _FieldOrderNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to FieldOrder values", s)
}

// FieldOrderValues returns all values of the enum
func FieldOrderValues() []FieldOrder {
	return _FieldOrderValues
}

// IsAFieldOrder returns "true" if the value is listed in the enum definition. "false" otherwise
func (i FieldOrder) IsAFieldOrder() bool {
	for _, v := range _FieldOrderValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for FieldOrder
func (i FieldOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for FieldOrder
func (i *FieldOrder) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FieldOrder should be a string, got %s", data)
	}

	var err error
	*i, err = FieldOrderString(s)
	return err
}
