package vm

import (
	"fmt"
)

var (
	// TRUE is shared boolean object that represents true
	TRUE *BooleanObject
	// FALSE is shared boolean object that represents false
	FALSE *BooleanObject
)

// BooleanObject represents boolean object in goby.
// It includes `true` and `FALSE` which represents logically true and false value.
// - `Boolean.new` is not supported.
type BooleanObject struct {
	*baseObj
	Value bool
}

func (vm *VM) initBoolClass() *RClass {
	b := vm.initializeClass(booleanClass, false)
	b.setBuiltInMethods(builtinBooleanInstanceMethods())
	b.class.setBuiltInMethods(builtInBooleanClassMethods())

	TRUE = &BooleanObject{Value: true, baseObj: &baseObj{class: b}}
	FALSE = &BooleanObject{Value: false, baseObj: &baseObj{class: b}}

	return b
}

func builtInBooleanClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.UnsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

func builtinBooleanInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Returns true if the receiver equals to the argument.
			//
			// ```Ruby
			// 1 == 1 # => true
			// 100 == 33 # => false
			// ```
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if receiver == args[0] {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns true if the receiver is not equals to the argument.
			//
			// ```Ruby
			// 4 != 2 # => true
			// 45 != 45 # => false
			// ```
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if receiver != args[0] {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Reverse the receiver.
			//
			// ```ruby
			// !true  # => false
			// !false # => true
			// ```
			// @return [Boolean]
			Name: "!",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					rightValue := receiver.(*BooleanObject).Value

					if rightValue {
						return FALSE
					}

					return TRUE
				}
			},
		},
		{
			// Returns true if both the receiver and the argument are true.
			//
			// ```ruby
			// 3 > 2 && 5 > 3  # => true
			// 3 > 2 && 5 > 10 # => false
			// ```
			// @return [Boolean]
			Name: "&&",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*BooleanObject).Value
					right, ok := args[0].(*BooleanObject)

					if !ok {
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, booleanClass, right.Class().Name)
						return err
					}

					rightValue := right.Value

					if leftValue && rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns true either if the receiver or argument is true.
			//
			// ```ruby
			// 3 > 2 || 5 > 3  # => true
			// 3 > 2 || 5 > 10 # => true
			// 2 > 3 || 5 > 10 # => false
			// ```
			// @return [Boolean]
			Name: "||",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*BooleanObject).Value
					right, ok := args[0].(*BooleanObject)

					if !ok {
						err := t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, booleanClass, right.Class().Name)
						return err
					}

					rightValue := right.Value

					if leftValue || rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString returns boolean object's value, which is either true or false.
func (b *BooleanObject) toString() string {
	return fmt.Sprintf("%t", b.Value)
}

// toJSON converts the receiver into JSON string.
func (b *BooleanObject) toJSON() string {
	return b.toString()
}

func (b *BooleanObject) equal(e *BooleanObject) bool {
	return b.Value == e.Value
}

func (b *BooleanObject) value() interface{} {
	return b.Value
}
