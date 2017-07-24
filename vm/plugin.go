package vm

import (
	"plugin"
	"reflect"
)

// PluginObject is a special type that contains a Go's plugin
type PluginObject struct {
	*baseObj
	fn     string
	plugin *plugin.Plugin
}

func (vm *VM) initPluginObject(fn string, p *plugin.Plugin) *PluginObject {
	c := vm.topLevelClass(pluginClass)
	return &PluginObject{fn: fn, plugin: p, baseObj: &baseObj{class: c, pseudoClass: c}}
}

func (vm *VM) initPluginClass() *RClass {
	pc := vm.initializeClass(pluginClass, false)
	pc.setBuiltInMethods(builtinPluginInstanceMethods())
	pc.class.setBuiltInMethods(builtinPluginClassMethods())

	return pc
}

// Only initialize file related methods after it's being required.
func builtinPluginClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}

// Only initialize file related methods after it's being required.
func builtinPluginInstanceMethods() []*BuiltInMethodObject {
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
					r := receiver.(*PluginObject)
					p := r.plugin
					f, err := p.Lookup(funcName)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					funcArgs, err := convertToGoFuncArgs(args)

					if err != nil {
						t.vm.initErrorObject(TypeError, err.Error())
					}

					result := reflect.ValueOf(reflect.ValueOf(f).Call(funcArgs)).Interface()

					return t.vm.initObjectFromGoType(unwrapGoFuncResult(result))
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString returns detailed infoof a array include elements it contains
func (p *PluginObject) toString() string {
	return "<Plugin: " + p.fn + ">"
}

// toJSON converts the receiver into JSON string.
func (p *PluginObject) toJSON() string {
	return p.toString()
}
