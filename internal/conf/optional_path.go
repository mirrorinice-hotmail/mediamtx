package conf

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/bluenviron/mediamtx/internal/conf/env"
	"github.com/bluenviron/mediamtx/internal/conf/jsonwrapper"
)

var optionalPathValuesType = func() reflect.Type {
	var fields []reflect.StructField
	rt := reflect.TypeOf(Path{})
	nf := rt.NumField()

	for i := 0; i < nf; i++ {
		f := rt.Field(i)
		j := f.Tag.Get("json")

		if j != "-" {
			if !strings.Contains(j, ",omitempty") {
				j += ",omitempty"
			}

			typ := f.Type
			if typ.Kind() != reflect.Pointer {
				typ = reflect.PointerTo(typ)
			}

			fields = append(fields, reflect.StructField{
				Name: f.Name,
				Type: typ,
				Tag:  reflect.StructTag(`json:"` + j + `"`),
			})
		}
	}

	return reflect.StructOf(fields)
}()

func newOptionalPathValues() interface{} {
	return reflect.New(optionalPathValuesType).Interface()
}

// OptionalPath is a Path whose values can all be optional.
type OptionalPath struct {
	ValuesOp interface{}
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *OptionalPath) UnmarshalJSON(b []byte) error {
	p.ValuesOp = newOptionalPathValues()
	return jsonwrapper.Unmarshal(b, p.ValuesOp)
}

// UnmarshalEnv implements env.Unmarshaler.
func (p *OptionalPath) UnmarshalEnv(prefix string, _ string) error {
	if p.ValuesOp == nil {
		p.ValuesOp = newOptionalPathValues()
	}
	return env.Load(prefix, p.ValuesOp)
}

// MarshalJSON implements json.Marshaler.
func (p *OptionalPath) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ValuesOp)
}
