package backend

import (
	"fmt"
	"github.com/zondax/keyringPoc/app/client"
)

func Backend(plugin string) {
	k, err := client.GetKeyring(plugin)
	if err != nil {
		panic(err)
	}
	defer k.Close()
	fmt.Println(k.Backend())
}
