package key

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/zondax/keyringPoc/app/client"
)

func output(r *keyring.Record) {
	add, err := r.GetAddress()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got key %s:\n", r.Name)
	fmt.Printf("Public key:\n\tType Url: %s\n\tValue: %v\n", r.PubKey.TypeUrl, r.PubKey.Value)
	fmt.Printf("Private key:\n\tType Url: %v\n\tValue: %v\n", r.GetLocal().PrivKey.TypeUrl, r.GetLocal().PrivKey.Value)
	fmt.Printf("Address:\n\t%s\n", add.String())

}

func Get(uid, plugin string) {
	k, err := client.GetKeyring(plugin)
	if err != nil {
		panic(err)
	}
	defer k.Close()
	key, err := k.Key(uid)
	if err != nil {
		panic(err)
	}

	output(key)
}
