package bindings

import (
	_ "embed"
)

const (
	Create3DeployedBytecode = "0x60806040908082526004918236101561001757600080fd5b600091823560e01c90816350f1c4641461027d575063cdcb760a1461003b57600080fd5b807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610279576024359067ffffffffffffffff82116102755736602383011215610275578184013593610090856103ba565b9361009d83519586610379565b8585526020958686019436602483830101116102755781839260248a930188378701015282513360601b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016878201908152833560148201529061012c81603484015b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282610379565b51902093846101396103f4565b888151910184f59073ffffffffffffffffffffffffffffffffffffffff9586831615610219579183929161016d849361042d565b98519134905af1903d15610213573d90610186826103ba565b9161019386519384610379565b8252873d92013e5b80610209575b156101ae57505191168152f35b8460649251917f08c379a0000000000000000000000000000000000000000000000000000000008352820152601560248201527f494e495449414c495a4154494f4e5f4641494c454400000000000000000000006044820152fd5b50833b15156101a1565b5061019b565b6064858a8851917f08c379a0000000000000000000000000000000000000000000000000000000008352820152601160248201527f4445504c4f594d454e545f4641494c45440000000000000000000000000000006044820152fd5b8280fd5b5080fd5b8284863461027957827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261027957359273ffffffffffffffffffffffffffffffffffffffff91828516850361032b575060609390931b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166020848101918252602435603486015293610323919061031b8160548101610100565b51902061042d565b915191168152f35b80fd5b6040810190811067ffffffffffffffff82111761034a57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761034a57604052565b67ffffffffffffffff811161034a57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b604051906104018261032e565b601082527f67363d3d37363d34f03d5260086018f3000000000000000000000000000000006020830152565b6104356103f4565b602081519101206040519060208201927fff0000000000000000000000000000000000000000000000000000000000000084523060601b602184015260358301526055820152605581526080810181811067ffffffffffffffff82111761034a577f010000000000000000000000000000000000000000000000000000000000000060b673ffffffffffffffffffffffffffffffffffffffff948360405284519020937fffffffffffffffffffffffffffffffffffffffff00000000000000000000000060a08201957fd694000000000000000000000000000000000000000000000000000000000000875260601b1660a28201520152601781526105398161032e565b519020169056fea26469706673582212203c3f094c14080645260151fc0973813e7e3e388692d8194b173368b56b8ac64864736f6c63430008170033"
)

//go:embed create3_storage_layout.json
var create3StorageLayoutJSON []byte

var Create3StorageLayout = mustGetStorageLayout(create3StorageLayoutJSON)