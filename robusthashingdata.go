package care

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func squashSequenceToString(val reflect.Value, builder *strings.Builder) (err error) {
	arrlen := val.Len()
	if arrlen < 1 {
		builder.WriteString("nil")
	} else {
		sb := strings.Builder{}
		items := make([]string, 0, arrlen)
		for i := 0; i < arrlen; i++ {
			err = robustHashingData(val.Index(i).Interface(), &sb)
			if err != nil {
				return
			}
			items = append(items, sb.String())
			sb.Reset()
		}
		sort.Strings(items)
		builder.WriteString(strings.Join(items, " "))
	}
	return
}

func squashMapToString(val reflect.Value, builder *strings.Builder) (err error) {
	mkeys := val.MapKeys()
	mlen := len(mkeys)
	if mlen < 1 {
		builder.WriteString("nil")
	} else {
		items := make([]string, 0, 2*mlen)
		sb := strings.Builder{}
		for _, mkey := range mkeys {
			mval := val.MapIndex(mkey)
			if err = robustHashingData(mkey.Interface(), &sb); err != nil {
				return
			}
			kitem := sb.String()
			sb.Reset()
			if err = robustHashingData(mval.Interface(), &sb); err != nil {
				return
			}
			vitem := sb.String()
			sb.Reset()
			items = append(items, kitem+" "+vitem)
		}
		sort.Strings(items)
		builder.WriteString(strings.Join(items, " "))
	}
	return
}

func squashStructToString(val reflect.Value, typ reflect.Type, builder *strings.Builder) (err error) {
	fields := val.NumField()
	items := make([]string, 0, fields)
	sb := strings.Builder{}
	for i := 0; i < fields; i++ {
		tf := typ.Field(i)
		if !tf.IsExported() {
			continue
		}
		if err = robustHashingData(val.Field(i).Interface(), &sb); err != nil {
			return
		}
		fname := tf.Name
		item := fname + " " + sb.String()
		sb.Reset()
		items = append(items, item)
	}
	sort.Strings(items)
	builder.WriteString(strings.Join(items, " "))
	return
}

func robustHashingData(val interface{}, builder *strings.Builder) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("Something goes wrong, ", r))
		}
	}()

	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return errors.New("The input value is not a valid value.")
	}

	typ := v.Type()

	builder.WriteString(typ.String() + ": ")

	switch v.Type().Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
		builder.WriteString(fmt.Sprint(v))
	case reflect.Array, reflect.Slice:
		err = squashSequenceToString(v, builder)
	case reflect.Map:
		err = squashMapToString(v, builder)
	case reflect.Pointer:
		if v.IsNil() {
			builder.WriteString("nil")
		} else {
			err = robustHashingData(v.Elem().Interface(), builder)
		}
	case reflect.Struct:
		err = squashStructToString(v, typ, builder)
	default:
		panic("The type \"" + typ.String() + "\" is not for processing.")
	}
	return
}
