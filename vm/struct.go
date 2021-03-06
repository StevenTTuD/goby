package vm

import (
	"fmt"
	"github.com/st0012/metago"
)

// StructObject ...
type StructObject struct {
	*baseObj
	data interface{}
}

func (vm *VM) initStructObject(d interface{}) *StructObject {
	return &StructObject{data: d, baseObj: &baseObj{class: vm.topLevelClass(structClass)}}
}

func (vm *VM) initStructClass() *RClass {
	sc := vm.initializeClass(structClass, false)
	sc.setBuiltInMethods(builtinStructClassMethods(), true)
	sc.setBuiltInMethods(builtinStructInstanceMethods(), false)
	vm.objectClass.setClassConstant(sc)
	return sc
}

// Only initialize file related methods after it's being required.
func builtinStructClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}

// Only initialize file related methods after it's being required.
func builtinStructInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "send",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, stringClass, args[0].Class().Name)
					}

					funcName := s.Value
					r := receiver.(*StructObject)

					funcArgs, err := convertToGoFuncArgs(args[1:])

					if err != nil {
						t.vm.initErrorObject(TypeError, err.Error())
					}

					result := metago.CallFunc(r.data, funcName, funcArgs...)
					return t.vm.initObjectFromGoType(result)
				}
			},
		},
	}
}

func convertToGoFuncArgs(args []Object) ([]interface{}, error) {
	funcArgs := []interface{}{}

	for _, arg := range args {
		v, ok := arg.(builtInType)

		if ok {
			if integer, ok := v.(*IntegerObject); ok {
				switch integer.flag {
				case integer64:
					funcArgs = append(funcArgs, int64(integer.Value))
					continue
				case integer32:
					funcArgs = append(funcArgs, int32(integer.Value))
					continue
				}
			}

			funcArgs = append(funcArgs, v.value())
		} else {
			err := fmt.Errorf("Can't pass %s type object when calling go function", arg.Class().Name)
			return nil, err
		}
	}

	return funcArgs, nil
}

// Polymorphic helper functions -----------------------------------------

func (s *StructObject) toString() string {
	return fmt.Sprintf("<Strcut: %p>", s)
}

func (s *StructObject) toJSON() string {
	return s.toString()
}
