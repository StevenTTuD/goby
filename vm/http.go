package vm

import (
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	httpRequestClass  *RClass
	httpResponseClass *RClass
)

func initHTTPClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	http := vm.initializeClass("HTTP", false)
	http.class.setBuiltInMethods(builtinHTTPClassMethods())

	initRequestClass(vm, http)
	initResponseClass(vm, http)

	net.setClassConstant(http)

	// Use Goby code to extend request and response classes.
	vm.execGobyLib("net/http/response.gb")
	vm.execGobyLib("net/http/request.gb")
}

func initRequestClass(vm *VM, hc *RClass) *RClass {
	requestClass := vm.initializeClass("Request", false)
	hc.setClassConstant(requestClass)
	builtinHTTPRequestInstanceMethods := []*BuiltInMethodObject{}

	requestClass.setBuiltInMethods(builtinHTTPRequestInstanceMethods)

	httpRequestClass = requestClass
	return requestClass
}

func initResponseClass(vm *VM, hc *RClass) *RClass {
	responseClass := vm.initializeClass("Response", false)
	hc.setClassConstant(responseClass)
	builtinHTTPResponseInstanceMethods := []*BuiltInMethodObject{}

	responseClass.setBuiltInMethods(builtinHTTPResponseInstanceMethods)

	httpResponseClass = responseClass
	return responseClass
}

func builtinHTTPClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "get",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					var path string

					domain := args[0].(*StringObject).Value

					if len(args) > 1 {
						path = args[1].(*StringObject).Value
					}

					if !strings.HasPrefix(path, "/") {
						path = "/" + path
					}

					resp, err := http.Get(domain + path)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		},
	}
}
