package stellar_test

import (
	"fmt"

	"github.com/soroworks/soroforge/internal/stellar"
	"github.com/stellar/go-stellar-sdk/xdr"
)

// ExampleContractIDFromSimulation recovers the contract ID a create-contract
// simulation reports, so a dry run can print the address a real deploy would
// produce before anything is submitted.
func ExampleContractIDFromSimulation() {
	// The simulated return value is the new contract's address, encoded as
	// base64 XDR. ContractIDFromSimulation decodes it back into a C... strkey.
	encoded, err := xdr.MarshalBase64(xdr.ScVal{
		Type: xdr.ScValTypeScvAddress,
		Address: &xdr.ScAddress{
			Type: xdr.ScAddressTypeScAddressTypeContract,
			ContractId: func() *xdr.ContractId {
				raw := make([]byte, 32)
				for i := range raw {
					raw[i] = byte(i + 1)
				}
				id := xdr.ContractId(raw)
				return &id
			}(),
		},
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	got, err := stellar.ContractIDFromSimulation(stellar.SimulateResult{
		Results: []stellar.HostFunctionResult{{ReturnValueXDR: encoded}},
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(got)

	// Output:
	// CAAQEAYEAUDAOCAJBIFQYDIOB4IBCEQTCQKRMFYYDENBWHA5DYPSBFLM
}

// ExampleCodeLedgerKey builds the ledger key used to look up uploaded WASM
// bytecode by its content hash.
func ExampleCodeLedgerKey() {
	key := stellar.CodeLedgerKey(stellar.WasmHash([]byte("soroforge")))

	fmt.Println(key.Type == xdr.LedgerEntryTypeContractCode)
	fmt.Println(stellar.HashHex(key.ContractCode.Hash))

	// Output:
	// true
	// dd2d49def30aef488ba257a346cec2b4b1f4d07a7bca033d4ae7dd716a8ff3e7
}

// ExampleInstanceLedgerKey builds the persistent ledger key for a deployed
// contract instance.
func ExampleInstanceLedgerKey() {
	key, err := stellar.InstanceLedgerKey("CAAQEAYEAUDAOCAJBIFQYDIOB4IBCEQTCQKRMFYYDENBWHA5DYPSBFLM")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(key.Type == xdr.LedgerEntryTypeContractData)
	fmt.Println(key.ContractData.Durability == xdr.ContractDataDurabilityPersistent)

	// Output:
	// true
	// true
}

// ExampleContractIDFromAsset derives the deterministic Stellar Asset Contract
// address for an asset on a specific network.
func ExampleContractIDFromAsset() {
	asset, err := stellar.ParseAsset("native")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	contractID, err := stellar.ContractIDFromAsset(asset, "Test SDF Network ; September 2015")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(len(contractID), contractID[0] == 'C')

	// Output:
	// 56 true
}

// ExampleContractIDFromDeployer derives the deterministic address a contract
// deployment will receive before any transaction is submitted.
func ExampleContractIDFromDeployer() {
	const deployer = "GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWHF"
	var salt [32]byte

	contractID, err := stellar.ContractIDFromDeployer(
		"Test SDF Network ; September 2015",
		deployer,
		salt,
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(len(contractID), contractID[0] == 'C')

	// Output:
	// 56 true
}
