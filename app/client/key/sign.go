package key

import (
	"fmt"
	"github.com/zondax/keyringPoc/app/client"
)

func Sign(uid, plugin, msg string) {
	k, err := client.GetKeyring(plugin)
	if err != nil {
		panic(err)
	}
	s, _, err := k.Sign(uid, []byte(msg))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(s))
}
