# This code has been obtained from:
# https://github.com/hukkin/cosmospy


import hashlib
import bech32
import mnemonic
import hdwallets
import ecdsa

from keyring.cosmos.crypto.keyring import v1 as cosmos_keyring
import betterproto.lib.google.protobuf as betterproto_lib_google_protobuf


PUBKEY_TYPE_URL = "/cosmos.crypto.secp256k1.PubKey"
PRIVKEY_TYPE_URL = "/cosmos.crypto.secp256k1.PrivKey"


def generate(seed: str, hd_path: str, bip39_passphrase: str) -> (bytes, bytes):
    # Generate priv
    seed_bytes = mnemonic.Mnemonic.to_seed(seed, bip39_passphrase)
    hd_wallet = hdwallets.BIP32.from_seed(seed_bytes)
    derived_privkey = hd_wallet.get_privkey_from_path(hd_path)
    privkey_obj = ecdsa.SigningKey.from_string(derived_privkey, curve=ecdsa.SECP256k1)
    # Generate pub
    pubkey_obj = privkey_obj.get_verifying_key()
    pubkey = pubkey_obj.to_string("compressed")
    s = hashlib.new("sha256", pubkey).digest()
    r = hashlib.new("ripemd160", s).digest()
    five_bit_r = bech32.convertbits(r, 8, 5)
    assert five_bit_r is not None, "Unsuccessful bech32.convertbits call"
    # Generate Addr
    address = bech32.bech32_encode("cosmos", five_bit_r)
    print(address)
    return derived_privkey, pubkey, address


def sign(record: cosmos_keyring.Record, msg: bytes) -> bytes:
    privkey = ecdsa.SigningKey.from_string(record.local.priv_key.value[2:], curve=ecdsa.SECP256k1)
    return privkey.sign_deterministic(msg, hashfunc=hashlib.sha256, sigencode=ecdsa.util.sigencode_string_canonize)


def new_record(uid: str, seed: str, hdpath: str, bip39_passphrase: str) -> cosmos_keyring.Record:
    priv, pub, addr = generate(seed, hdpath, bip39_passphrase)
    return cosmos_keyring.Record(
        name=uid,
        pub_key=betterproto_lib_google_protobuf.Any(
            type_url=PUBKEY_TYPE_URL,
            value=pub,
        ),
        local=cosmos_keyring.RecordLocal(
            priv_key=betterproto_lib_google_protobuf.Any(
                type_url=PRIVKEY_TYPE_URL,
                value=priv
            )
        )
    )


def new_record_offline(uid: str, pub_key: betterproto_lib_google_protobuf.Any) -> cosmos_keyring.Record:
    return cosmos_keyring.Record(
        name=uid,
        pub_key=betterproto_lib_google_protobuf.Any(
            type_url=pub_key.type_url,
            value=pub_key.value,
        ),
        offline=cosmos_keyring.RecordOffline()
    )
