import asyncio
import sys

from keyring.keyring import v1 as keyring
from keyring.cosmos.crypto.keyring import v1 as cosmos_keyring

from grpclib.server import Server
import secp256k1


def marshall_record(record: cosmos_keyring.Record):
    record.pub_key.value = bytes("\n!", 'utf-8') + record.pub_key.value
    record.local.priv_key.value = bytes("\n ", 'utf-8') + record.local.priv_key.value
    return record


class PymemService(keyring.KeyringServiceBase):
    db = dict()

    async def backend(self, backend_request: "keyring.BackendRequest") -> "keyring.BackendResponse":
        return keyring.BackendResponse(backend="pymem")

    async def key(self, key_request: "keyring.KeyRequest") -> "keyring.KeyResponse":
        key = self.db[key_request.uid]
        return keyring.KeyResponse(key=key)

    async def new_account(
            self, new_account_request: "keyring.NewAccountRequest"
    ) -> "keyring.NewAccountResponse":
        record = secp256k1.new_record(new_account_request.uid,
                                      new_account_request.mnemonic,
                                      new_account_request.hdpath,
                                      new_account_request.bip39_passphrase)
        record = marshall_record(record)
        self.db[new_account_request.uid] = bytes(record)
        return keyring.NewAccountResponse(record=record)

    async def sign(self, new_sign_request: "keyring.NewSignRequest") -> "keyring.NewSignResponse":
        record = cosmos_keyring.Record.FromString(self.db[new_sign_request.uid])
        return keyring.NewSignResponse(
            msg=secp256k1.sign(record, new_sign_request.msg)
        )


async def serve():
    host = "127.0.0.1"
    port = 10009
    server = Server([PymemService()])
    await server.start(host, port)
    print("1|1|tcp|127.0.0.1:10009|grpc")
    sys.stdout.flush()
    await server.wait_closed()


if __name__ == "__main__":
    s = serve()
    asyncio.run(s)
