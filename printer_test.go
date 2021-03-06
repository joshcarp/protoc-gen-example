package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/gogo/protobuf/proto"
	plugin_go "github.com/gogo/protobuf/protoc-gen-gogo/plugin"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/syslutil"

	pgs "github.com/lyft/protoc-gen-star"
	"github.com/spf13/afero"
)

func TestPrinting(t *testing.T) {
	directory := "./tests/test/"
	_, fs := syslutil.WriteToMemOverlayFs(directory)

	GeneratorResponse, err := ConvertSyslToProto(directory + "code_generator_request.pb.bin")
	assert.NoError(t, err)

	golden, err := afero.ReadFile(fs, *GeneratorResponse.File[0].Name)
	assert.NoError(t, err)

	assert.Equal(t, *GeneratorResponse.File[0].Content, string(golden))
}

func TestSimple(t *testing.T) {
	directory := "./tests/simple/"
	_, fs := syslutil.WriteToMemOverlayFs(directory)
	GeneratorResponse, err := ConvertSyslToProto(directory + "code_generator_request.pb.bin")
	assert.NoError(t, err)

	golden, err := afero.ReadFile(fs, *GeneratorResponse.File[0].Name)
	assert.NoError(t, err)

	assert.Equal(t, *GeneratorResponse.File[0].Content, string(golden))
}

func ConvertSyslToProto(filename string) (*plugin_go.CodeGeneratorResponse, error) {
	req, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	res := &bytes.Buffer{}
	pgs.Init(
		pgs.ProtocInput(req),  // use the pre-generated request
		pgs.ProtocOutput(res), // capture CodeGeneratorResponse
	).RegisterModule(
		SyslPrinter(),
	).Render()
	response := plugin_go.CodeGeneratorResponse{}
	err = proto.Unmarshal(res.Bytes(), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil

}
