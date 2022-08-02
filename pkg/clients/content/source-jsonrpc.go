// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package content

import (
	"context"
	"encoding/json"
	goUUID "github.com/google/uuid"
)

type ClientSource struct {
	*ClientJsonRPC
}

type retSourceMigration func(err error)

func (cli *ClientSource) ReqMigration(ret retSourceMigration, attributes ...map[string]interface{}) (request baseJsonRPC) {

	request = baseJsonRPC{
		Method:  "source.migration",
		Params:  requestSourceMigration{Attributes: attributes},
		Version: Version,
	}
	var err error
	var response responseSourceMigration

	if ret != nil {
		request.retHandler = func(jsonrpcResponse baseJsonRPC) {
			if jsonrpcResponse.Error != nil {
				err = cli.errorDecoder(jsonrpcResponse.Error)
				ret(err)
				return
			}
			err = json.Unmarshal(jsonrpcResponse.Result, &response)
			ret(err)
		}
		request.ID = []byte("\"" + goUUID.New().String() + "\"")
	}
	return
}

func (cli *ClientSource) Migration(ctx context.Context, attributes ...map[string]interface{}) (err error) {

	retHandler := func(_err error) {
		err = _err
	}
	if blockErr := cli.Batch(ctx, cli.ReqMigration(retHandler, attributes...)); blockErr != nil {
		err = blockErr
		return
	}
	return
}