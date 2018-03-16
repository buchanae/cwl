package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

// loader helps deal with type coercion while loading a
// CWL document, for example, dealing with the fact that
// "inputs" might be a scalar, a map, or a list.
//
// loader tries to detect obvious type coercions via reflect.
// non-obvious type coercions must have a registered
// handler to do the work.
type loader struct {
	base     string
	resolver Resolver
}

// load is given a YAML node and a destination type,
// e.g. yamlast.Mapping -> cwl.WorkflowInput.
//
// load() panics if `t` is not a pointer.
// load() panics if given an unknown YAML node type (such as Alias)
func (l *loader) load(n node, t interface{}) error {

	// only pointers can be set to new values by the loader.
	if reflect.TypeOf(t).Kind() != reflect.Ptr {
		return fmt.Errorf("load() must be called with a pointer")
	}

	// get the reflected type of the loader in order to look up
	// handler methods, e.g. loader.MappingToWorkflowInput()
	loaderTyp := reflect.TypeOf(l)
	loaderVal := reflect.ValueOf(l)

	// reflect the type and value of the destination.
	typ := reflect.TypeOf(t).Elem()
	val := reflect.ValueOf(t).Elem()

	// get string version of the yaml node type
	// for building the handler name.
	nodeKind := "unknown"
	switch n.Kind {
	case yamlast.MappingNode:
		nodeKind = "Mapping"
	case yamlast.SequenceNode:
		nodeKind = "Seq"
	case yamlast.ScalarNode:
		nodeKind = "Scalar"
	default:
		panic("unknown node kind")
	}

	// describes the type conversion being requested,
	// in order to look up a registered handler.
	typename := strings.Title(typ.Name())
	if typ.Kind() == reflect.Slice {
		typename = strings.Title(typ.Elem().Name())
		typename += "Slice"
	}
	if typ.Kind() == reflect.Map {
		typename = strings.Title(typ.Elem().Name())
		typename += "Map"
	}
	handlerName := nodeKind + "To" + typename

	// look for a handler. if found, use it.
	if _, ok := loaderTyp.MethodByName(handlerName); ok {
		m := loaderVal.MethodByName(handlerName)
		nval := reflect.ValueOf(n)
		outv := m.Call([]reflect.Value{nval})
		errv := outv[1]
		if !errv.IsNil() {
			return errv.Interface().(error)
		}
		resv := outv[0]
		val.Set(resv)
		return nil
	}

	switch {
	// Try to handle obvious scalar conversions automatically.
	case n.Kind == yamlast.ScalarNode:
		vt := reflect.TypeOf(n.Value)

		if vt.AssignableTo(typ) {
			val.Set(reflect.ValueOf(n.Value))
			return nil
		} else if vt.ConvertibleTo(typ) {
			val.Set(reflect.ValueOf(n.Value).Convert(typ))
			return nil
		} else {
			err := coerceSet(t, n.Value)
			if err == nil {
				return nil
			}
		}

	// Try to automatically load a YAML mapping into a struct.
	case typ.Kind() == reflect.Struct && n.Kind == yamlast.MappingNode:
		return l.loadMappingToStruct(n, t)

		// Try to automatically load a YAML sequence into a slice type,
		// without a defined handler.
	case typ.Kind() == reflect.Slice && n.Kind == yamlast.SequenceNode:
		for _, c := range n.Children {
			el := typ.Elem()
			if el.Kind() == reflect.Ptr {
				el = el.Elem()
			}

			item := reflect.New(el)
			err := l.load(c, item.Interface())
			if err != nil {
				return err
			}

			if typ.Elem().Kind() == reflect.Ptr {
				val.Set(reflect.Append(val, item))
			} else {
				val.Set(reflect.Append(val, item.Elem()))
			}
		}
		return nil
	}

	// No handler found.
	//debug(handlerName)
	return fmt.Errorf("unhandled type at line %d, col %d", n.Line+1, n.Column+1)
}

// loadMappingToStruct essentially unmarshals a YAML mapping
// into a Go struct.
//
// "n" must be a mapping node.
// "t" must be a pointer to a struct.
func (l *loader) loadMappingToStruct(n node, t interface{}) error {

	if n.Kind != yamlast.MappingNode {
		panic("expected mapping node")
	}
	if len(n.Children)%2 != 0 {
		panic("expected even number of children in mapping")
	}

	typ := reflect.TypeOf(t).Elem()
	val := reflect.ValueOf(t).Elem()
	// track which fields have been set in order to raise an error
	// when a field exists twice.
	already := map[string]bool{}

	for i := 0; i < len(n.Children)-1; i += 2 {
		k := n.Children[i]
		v := n.Children[i+1]
		name := strings.ToLower(k.Value)

		if _, ok := already[name]; ok {
			return fmt.Errorf("duplicate field found while loading mapping")
		}
		already[name] = true

		// Find a matching field in the target struct.
		// Names are case insensitive.
		var field reflect.StructField
		var found bool
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)

			n := f.Name
			if alt, ok := f.Tag.Lookup("json"); ok {
				sp := strings.Split(alt, ",")
				n = sp[0]
			}

			if strings.ToLower(n) == name {
				field = f
				found = true
				break
			}
		}

		if !found {
			continue
		}

		fv := val.FieldByIndex(field.Index)

		if !fv.CanSet() {
			continue
		}

		var val reflect.Value
		if field.Type.Kind() == reflect.Ptr {
			val = reflect.New(field.Type.Elem())
			fv.Set(val)
		} else {
			val = fv.Addr()
		}

		err := l.load(v, val.Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

// coerceSet attempts to coerce "val" to the type of "dest".
// If coercion succeeds, "dest" is set to the coerced value of "val".
// coerceSet panics if "dest" is not a pointer.
func coerceSet(dest interface{}, val interface{}) error {
	var casted interface{}
	var err error

	switch dest.(type) {
	case *int:
		casted, err = cast.ToIntE(val)
	case *int64:
		casted, err = cast.ToInt64E(val)
	case *int32:
		casted, err = cast.ToInt32E(val)
	case *float32:
		casted, err = cast.ToFloat32E(val)
	case *float64:
		casted, err = cast.ToFloat64E(val)
	case *bool:
		casted, err = cast.ToBoolE(val)
	case *string:
		casted, err = cast.ToStringE(val)
	case *[]string:
		casted, err = cast.ToStringSliceE(val)
	default:
		return fmt.Errorf("unknown dest type: %s", dest)
	}

	if err != nil {
		return fmt.Errorf("error casting: %s", err)
	}

	reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(casted))
	return nil
}
