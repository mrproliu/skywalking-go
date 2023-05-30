package core

import (
	"reflect"
	"unsafe"
)

type TracerTools struct {
}

func NewTracerTools() *TracerTools {
	return &TracerTools{}
}

type ReflectFieldFilter struct {
	name          string
	interfaceType interface{}
	typeVal       interface{}
}

func (r *ReflectFieldFilter) SetName(val string) {
	r.name = val
}

func (r *ReflectFieldFilter) SetInterfaceType(val interface{}) {
	r.interfaceType = val
}

func (r *ReflectFieldFilter) SetType(val interface{}) {
	r.typeVal = val
}

type ReflectFieldFilterOpts interface {
	Apply(interface{})
}

func (t *TracerTools) ReflectGetValueByType(instance interface{}, filterOpts []interface{}) interface{} {
	instanceVal := reflect.ValueOf(instance)
	if instanceVal.Kind() == reflect.Ptr && instanceVal.Elem().Kind() == reflect.Struct {
		instanceVal = instanceVal.Elem()
	}
	filter := &ReflectFieldFilter{}
	for _, opt := range filterOpts {
		if f, ok := opt.(ReflectFieldFilterOpts); ok {
			f.Apply(filter)
		}
	}
	for i := 0; i < instanceVal.NumField(); i++ {
		field := instanceVal.Field(i)
		fieldType := instanceVal.Type().Field(i)

		if t.checkFieldSupport(field, &fieldType, filter) {
			// for getting the export field value
			// 1. get a pointer to the field, then convert it to a 'generic' pointer
			// 2. convert the pointer back to an interface{}
			fieldPtr := unsafe.Pointer(field.UnsafeAddr())
			return reflect.NewAt(field.Type(), fieldPtr).Elem().Interface()
		}
	}
	return nil
}

func (t *TracerTools) checkFieldSupport(field reflect.Value, instanceField *reflect.StructField, filter *ReflectFieldFilter) bool {
	if filter.name != "" {
		if instanceField.Name != filter.name {
			return false
		}
	}
	if filter.interfaceType != nil {
		interfaceType := reflect.TypeOf(filter.interfaceType).Elem()
		if !field.Type().Implements(interfaceType) {
			return false
		}
	}
	if filter.typeVal != nil {
		if field.Interface() != filter.typeVal {
			return false
		}
	}
	return true
}
