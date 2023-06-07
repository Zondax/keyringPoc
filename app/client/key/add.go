package key

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/zondax/keyringPoc/app/client"
)

func outputKey(r *keyring.Record, mnemonic string) {
	fmt.Printf("Key %s created:\n", r.Name)
	fmt.Printf("Public key:\n\tType Url: %s\n\tValue: %v\n", r.PubKey.TypeUrl, r.PubKey.Value)
	fmt.Printf("Private key:\n\tType Url: %v\n\tValue: %v\n", r.GetLocal().PrivKey.TypeUrl, r.GetLocal().PrivKey.Value)
	fmt.Println("mnemonic:")
	fmt.Printf("\t%s", mnemonic)
}

func Add(uid, plugin, mnemonic string) {
	k, err := client.GetKeyring(plugin)
	if err != nil {
		panic(err)
	}
	if mnemonic == "" {
		mnemonic, err = client.NewMnemonic()
		if err != nil {
			panic(err)
		}
	}
	mnemonic = "spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police"
	r, err := k.NewAccount(uid, mnemonic, "", "m/44'/118'/0'/0/0", hd.Secp256k1)
	if err != nil {
		panic(err)
	}
	outputKey(r, mnemonic)
}
