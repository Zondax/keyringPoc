# Keyring POC

This repository contains a proof of concept (POC) for a keyring implementation that leverages HashiCorp plugins over
gRPC. The keyring is designed to integrate with the Cosmos SDK and fulfill the current keyring interface of the SDK.
The POC focuses on implementing the methods for key creation, retrieval, and signing.


## Overview
The main objective of this POC is to explore the viability of using HashiCorp plugins over gRPC within a keyring
context. It aims to demonstrate how the keyring can interact with external plugins written in both Go and Python,
enabling seamless integration with the Cosmos SDK.

This POC only implements the following methods:
 * key creation
 * key retrieving
 * signing of messages. 

The gRPC service is defined as follows:
```protobuf
service KeyringService {
    rpc Backend(BackendRequest) returns (BackendResponse);
    rpc Key(KeyRequest) returns (KeyResponse);
    rpc NewAccount(NewAccountRequest) returns (NewAccountResponse);
    rpc Sign(NewSignRequest) returns (NewSignResponse);
}
```
you can check all the [proto messages](proto/keyringPoc/keyring/v1/request.proto) under the proto dir.


The [keystore](keyring/keyStore/keyStore.go) fulfils the [cosmos-sdk keyring](https://github.com/cosmos/cosmos-sdk/blob/v0.47.3-rc.0/crypto/keyring/keyring.go)
interface but only the following methods have been implemented:
```go
    Backend() string
    Key(uid string) (*Record, error)
    NewAccount(uid, mnemonic, bip39Passphrase, hdPath string, algo SignatureAlgo) (*Record, error)
    Sign(uid string, msg []byte, signMode signing.SignMode) ([]byte, types.PubKey, error)
```


## Proto
To compile the protobuf files, run:

`make proto-gen`

This will generate the necessary go and python files.


## Use
Before using the POC, you need to compile the Go plugin and download the Python dependencies.

### Go plugin
To compile go plugin, run:
 * `make build-go-plugin-file`

### Python plugin
To use the Python plugin, ensure that Python 3 is installed on your system along with [virtualenv](https://virtualenv.pypa.io/en/latest/).
Once virtualenv is set up, execute the following commands:

* `make python-setup`
* `source .venv/bin/activate`

this will download all the dependencies and activate the environment.

### App
A CLI app is available to facilitate the use of this keyring. To build the app, run:
* `make build-app`

This will create a binary named `app` under build directory, which you can use.

To check if everything is working smoothly:
 * `./build/app backend --plugin goFile`
 * `./build/app backend --plugin pyFile`

For more information on using the app, run `./build/app --help`

## Issues

### Record over gRCP
In the cosmos-sdk `google.protobuf.Any` is replaced with the following struct:
```go
type Any struct {
	// A URL/resource name that uniquely identifies the type of the serialized
	// protocol buffer message.
	TypeUrl string `protobuf:"bytes,1,opt,name=type_url,json=typeUrl,proto3"
                    json:"type_url,omitempty"`
	// Must be a valid serialized protocol buffer of the above specified type.
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
	cachedValue interface{}
	compat *anyCompat
}
```
This struct relies on its cachedValue field to retrieve the address and public key in the case of a Record.
When returning a Record over gRPC, the cachedValue is nil since it's not an exported attribute. Therefore, the
deserialization of the record must be performed within the SDK and cannot be accomplished within the plugins.

### SignatureAlgo as parameters
In the current implementation of the cosmos-sdk keyring, the following methods require a SignatureAlgo as a parameter:
```go
    NewMnemonic(uid string, language Language, hdPath, bip39Passphrase string, algo SignatureAlgo) (*Record, string, error)
    NewAccount(uid, mnemonic, bip39Passphrase, hdPath string, algo SignatureAlgo) (*Record, error)
    SaveLedgerKey(uid string, algo SignatureAlgo, hrp string, coinType, account, index uint32) (*Record, error)

```
The SignatureAlgo interface is defined as follows:
```go
// SignatureAlgo defines the interface for a keyring supported algorithm.
type SignatureAlgo interface {
	Name() hd.PubKeyType
	Derive() hd.DeriveFn
	Generate() hd.GenerateFn
}
```
However, since the SignatureAlgo information cannot be sent over gRPC, and the plugins are currently hardcoded to use
secp256k1, a potential solution would be to include the Name() of the algorithm as part of the request. In this
scenario, it would be reasonable to refactor the SDK to pass the algorithm as a string in those methods. Then, a lookup
can be performed on a map to obtain the corresponding SignatureAlgo based on the provided string.