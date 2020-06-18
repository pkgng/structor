package structor

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/robertkrimen/otto"
)

type Structor struct {
	target     interface{}
	langis     *otto.Otto
	components map[string]interface{}
	copy       bool
}

func NewStructor(target interface{}) *Structor {
	t := &Structor{}
	t.copy = false
	t.langis = otto.New()
	t.target = target

	t.langis.Set("self", target)
	return t
}

func (s *Structor) Set(name string, value interface{}) *Structor {
	s.components = make(map[string]interface{})
	s.components[name] = value
	s.langis.Set(name, value)
	return s
}

func (s *Structor) CopyByName() *Structor {
	s.copy = true
	return s
}

func execute(langis *otto.Otto, script string) interface{} {
	value, _ := langis.Run(script)
	r, e := value.Export()
	if e != nil {
		fmt.Println(e)
	}
	return r
}

func (s *Structor) Construct() error {
	if s.copy {
		for _, com := range s.components {
			copit(s.target, com)
		}
	}
	return s.calc()
}

// calculate fields
func (s *Structor) calc() (err error) {
	var (
		to = indirect(reflect.ValueOf(s.target))
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
		script := field.Tag.Get("structor")
		if script == "" {
			continue
		}

		r := execute(s.langis, script)

		// fmt.Printf("%s -> %v\n", script, r)
		if toField := to.FieldByName(field.Name); toField.IsValid() {
			if toField.CanSet() {
				if !set(toField, indirect(reflect.ValueOf(r))) {
					// fmt.Println("--------------- deep copy")
					// if err := copit(toField.Addr().Interface(), r); err != nil {
					// 	return err
					// }
				}
			}
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
		return
	}

	if from.IsValid() {
		toTypeFields := deepFields(toType)
		for _, field := range toTypeFields {
			name := field.Name

			if fromField := from.FieldByName(name); fromField.IsValid() {
				if toField := to.FieldByName(name); toField.IsValid() {
					if toField.CanSet() {
						if !set(toField, fromField) {
							if err := copit(toField.Addr().Interface(), fromField.Interface()); err != nil {
								return err
							}
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
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
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
