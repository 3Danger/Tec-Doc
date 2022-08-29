// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package suppliers

import (
	"context"
	"encoding/json"
	goUUID "github.com/google/uuid"
	"github.com/satori/go.uuid"
)

type ClientSuppliers struct {
	*ClientJsonRPC
}

type retSuppliersGetOldSupplierID func(oldSupplierID int, err error)

func (cli *ClientSuppliers) ReqGetOldSupplierID(ret retSuppliersGetOldSupplierID, supplierID uuid.UUID) (request baseJsonRPC) {

	request = baseJsonRPC{
		Method:  "suppliers.getoldsupplierid",
		Params:  requestSuppliersGetOldSupplierID{SupplierID: supplierID},
		Version: Version,
	}
	var err error
	var response responseSuppliersGetOldSupplierID

	if ret != nil {
		request.retHandler = func(jsonrpcResponse baseJsonRPC) {
			if jsonrpcResponse.Error != nil {
				err = cli.errorDecoder(jsonrpcResponse.Error)
				ret(response.OldSupplierID, err)
				return
			}
			err = json.Unmarshal(jsonrpcResponse.Result, &response)
			ret(response.OldSupplierID, err)
		}
		request.ID = []byte("\"" + goUUID.New().String() + "\"")
	}
	return
}

func (cli *ClientSuppliers) GetOldSupplierID(ctx context.Context, supplierID uuid.UUID) (oldSupplierID int, err error) {

	retHandler := func(_oldSupplierID int, _err error) {
		oldSupplierID = _oldSupplierID
		err = _err
	}
	if blockErr := cli.Batch(ctx, cli.ReqGetOldSupplierID(retHandler, supplierID)); blockErr != nil {
		err = blockErr
		return
	}
	return
}