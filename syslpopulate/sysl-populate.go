package syslpopulate

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

var TypeMapping = map[string]sysl.Type_Primitive{
	"TYPE_BYTES":  sysl.Type_BYTES,
	"TYPE_INT32":  sysl.Type_INT,
	"TYPE_STRING": sysl.Type_STRING,
	"TYPE_BOOL":   sysl.Type_BOOL,
}

// NewApplication Initialises a Sysl application
func NewApplication(appName string) *sysl.Application {
	return &sysl.Application{
		Name:      NewAppName(appName),
		Endpoints: map[string]*sysl.Endpoint{},
		Types:     map[string]*sysl.Type{},
	}
}

// NewEndpoint Initialises a Sysl Endpoint
func NewEndpoint(name string) *sysl.Endpoint {
	return &sysl.Endpoint{Name: name}
}

// NewParameter Initialises a Sysl Parameter input
func NewParameter(name, application string) *sysl.Param {
	return &sysl.Param{
		Name: "input",
		Type: NewType(name, application),
	}
}

// NewType Initialises a Sysl type from string
func NewType(name, application string) *sysl.Type {
	if fieldType, ok := TypeMapping[name]; ok {
		return SyslPrimitive(fieldType)
	}
	return SyslStruct(name, application)
}

// NewReturn Initialises a return statement and wraps it in a sysl statement
// payloads will be concatenated and seperated by dots "."
func NewReturn(payload ...string) *sysl.Statement {
	return &sysl.Statement{Stmt: &sysl.Statement_Ret{Ret: &sysl.Return{
		Payload: "ok <: " + strings.Join(payload, ".")}}}
}

// NewCall Initialises a call statement and wraps it in a sysl statement
func NewCall(app, endpoint string) *sysl.Statement {
	return &sysl.Statement{Stmt: &sysl.Statement_Call{
		Call: &sysl.Call{
			Target:   NewAppName(app),
			Endpoint: endpoint}}}
}

func NewAppName(name string) *sysl.AppName {
	return &sysl.AppName{Part: []string{name}}
}

// SyslPrimitive converts a string to a sysl primitive type (parameter must be in sysl type)
func SyslPrimitive(fieldType sysl.Type_Primitive) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Primitive_{
			Primitive: fieldType,
		},
	}
}

// SyslStruct converts a string to a sysl struct type
func SyslStruct(fieldType, application string) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Ref: &sysl.Scope{
					Appname: NewAppName(application),
					Path:    []string{fieldType},
				},
			},
		},
	}
}
