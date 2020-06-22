package structor

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/robertkrimen/otto"
)

type BaseStructor struct {
	target interface{}
	langis *otto.Otto
	copy   bool
}

type Structor struct {
	components map[string]interface{}
}

func New() *Structor {
	t := &Structor{}
	t.components = make(map[string]interface{})
	return t
}

func (s *Structor) Set(name string, value interface{}) *Structor {
	s.components[name] = value
	return s
}

func (s *Structor) getStructorBase(t interface{}) *BaseStructor {
	target := indirect(reflect.ValueOf(t))
	tt, haveBase := indirectType(target.Type()).FieldByName("BaseStructor")

	if !haveBase {
		return nil
	}

	if indirectType(tt.Type) == reflect.TypeOf(BaseStructor{}) {
		base := target.FieldByName("BaseStructor").Addr().Interface().(*BaseStructor)
		tags := tt.Tag.Get("structor")
		base.langis = otto.New()
		base.copy = false

		if strings.Contains(tags, "CopyByName") {
			base.copy = true
		}

		for _, tag := range strings.Split(tags, ",") {
			if tag != "CopyByName" {
				if base.copy {
					copit(t, s.components[tag])
				}

				base.langis.Set(tag, s.components[tag])
			}
		}
		base.langis.Set("self", t)

		return base
	}
	return nil
}

func execute(langis *otto.Otto, script string) (interface{}, error) {
	value, e := langis.Run(script)
	r, _ := value.Export()
	if e != nil {
		return nil, e
	}
	return r, nil
}

func (s *Structor) Construct(target interface{}) error {
	base := s.getStructorBase(target)
	if base == nil {
		return errors.New("no structor Base inSide")
	}

	calc(s, base, target)

	return nil
}

// calculate fields
func calc(root *Structor, base *BaseStructor, target interface{}) (err error) {
	var (
		to = indirect(reflect.ValueOf(target))
	)

	if !to.CanAddr() || !to.IsValid() {
		return errors.New("input value is unaddressable or inValid")
	}

	toType := indirectType(to.Type())

	// ---------------------  struct -> struct only
	if toType.Kind() != reflect.Struct {
		return errors.New("only struct accept")
	}

	toTypeFields := deepFields(toType)
	for _, field := range toTypeFields {
		if field.Name == "BaseStructor" {
			continue
		}

		toField := to.FieldByName(field.Name)
		if toField.Kind() == reflect.Struct {
			// fmt.Printf("deep in %s @ %s\n", field.Name, field.Tag.Get("structor"))
			base2 := root.getStructorBase(toField.Addr().Interface())
			if base2 == nil {
				calc(root, base, toField.Addr().Interface())
			} else {
				calc(root, base2, toField.Addr().Interface())
			}

			continue
		}

		script := field.Tag.Get("structor")
		if script == "" || !toField.IsValid() || !toField.CanSet() {
			continue
		}

		if r, err := execute(base.langis, script); err == nil {
			// fmt.Printf("%s -> %v\n", script, r)
			set(toField, indirect(reflect.ValueOf(r)))
		}
	}

	return
}

// Copy copy things
func copit(toValue interface{}, fromValue interface{}) (err error) {
	var (
		from = indirect(reflect.ValueOf(fromValue))
		to   = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() || !from.IsValid() {
		return errors.New("copy to value is unaddressable")
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	// Just set it if possible to assign
	// And need to do copy anyway if the type is struct
	if fromType.Kind() != reflect.Struct && from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	// ---------------------  struct -> struct only
	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return errors.New("only struct -> struct allowed")
	}

	if from.IsValid() {
		toTypeFields := deepFields(toType)
		for _, field := range toTypeFields {
			name := field.Name

			if fromField := from.FieldByName(name); fromField.IsValid() {
				if toField := to.FieldByName(name); toField.IsValid() {
					if toField.CanSet() {
						if !set(toField, fromField) {
							copit(toField.Addr().Interface(), fromField.Interface())
						}
					}
				}
			}
		}
	}

	return
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			fields = append(fields, v)
		}
	}

	return fields
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}
