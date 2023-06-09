# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: keyringPoc/keyring/v1/request.proto, keyringPoc/keyring/v1/service.proto
# plugin: python-betterproto
from dataclasses import dataclass
from typing import (
    TYPE_CHECKING,
    Dict,
    Optional,
)

import betterproto
import betterproto.lib.google.protobuf as betterproto_lib_google_protobuf
import grpclib
from betterproto.grpc.grpclib_server import ServiceBase

from ...cosmos.crypto.keyring import v1 as __cosmos_crypto_keyring_v1__
from ...cosmos.tx.signing import v1beta1 as __cosmos_tx_signing_v1_beta1__


if TYPE_CHECKING:
    import grpclib.server
    from betterproto.grpc.grpclib_client import MetadataLike
    from grpclib.metadata import Deadline


@dataclass(eq=False, repr=False)
class Empty(betterproto.Message):
    pass


@dataclass(eq=False, repr=False)
class BackendRequest(betterproto.Message):
    pass


@dataclass(eq=False, repr=False)
class BackendResponse(betterproto.Message):
    backend: str = betterproto.string_field(1)


@dataclass(eq=False, repr=False)
class KeyRequest(betterproto.Message):
    uid: str = betterproto.string_field(1)


@dataclass(eq=False, repr=False)
class KeyResponse(betterproto.Message):
    key: bytes = betterproto.bytes_field(1)


@dataclass(eq=False, repr=False)
class NewAccountRequest(betterproto.Message):
    uid: str = betterproto.string_field(1)
    mnemonic: str = betterproto.string_field(2)
    bip39_passphrase: str = betterproto.string_field(3)
    hdpath: str = betterproto.string_field(4)


@dataclass(eq=False, repr=False)
class NewAccountResponse(betterproto.Message):
    record: "__cosmos_crypto_keyring_v1__.Record" = betterproto.message_field(1)


@dataclass(eq=False, repr=False)
class SignRequest(betterproto.Message):
    uid: str = betterproto.string_field(1)
    msg: bytes = betterproto.bytes_field(2)
    sign_mode: "__cosmos_tx_signing_v1_beta1__.SignMode" = betterproto.enum_field(3)


@dataclass(eq=False, repr=False)
class SignResponse(betterproto.Message):
    msg: bytes = betterproto.bytes_field(1)
    pub_key: "betterproto_lib_google_protobuf.Any" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class SaveOfflineRequest(betterproto.Message):
    uid: str = betterproto.string_field(1)
    pub_key: "betterproto_lib_google_protobuf.Any" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class SaveOfflineResponse(betterproto.Message):
    record: bytes = betterproto.bytes_field(1)


class KeyringServiceStub(betterproto.ServiceStub):
    async def backend(
        self,
        backend_request: "BackendRequest",
        *,
        timeout: Optional[float] = None,
        deadline: Optional["Deadline"] = None,
        metadata: Optional["MetadataLike"] = None
    ) -> "BackendResponse":
        return await self._unary_unary(
            "/keyring.v1.KeyringService/Backend",
            backend_request,
            BackendResponse,
            timeout=timeout,
            deadline=deadline,
            metadata=metadata,
        )

    async def key(
        self,
        key_request: "KeyRequest",
        *,
        timeout: Optional[float] = None,
        deadline: Optional["Deadline"] = None,
        metadata: Optional["MetadataLike"] = None
    ) -> "KeyResponse":
        return await self._unary_unary(
            "/keyring.v1.KeyringService/Key",
            key_request,
            KeyResponse,
            timeout=timeout,
            deadline=deadline,
            metadata=metadata,
        )

    async def new_account(
        self,
        new_account_request: "NewAccountRequest",
        *,
        timeout: Optional[float] = None,
        deadline: Optional["Deadline"] = None,
        metadata: Optional["MetadataLike"] = None
    ) -> "NewAccountResponse":
        return await self._unary_unary(
            "/keyring.v1.KeyringService/NewAccount",
            new_account_request,
            NewAccountResponse,
            timeout=timeout,
            deadline=deadline,
            metadata=metadata,
        )

    async def sign(
        self,
        sign_request: "SignRequest",
        *,
        timeout: Optional[float] = None,
        deadline: Optional["Deadline"] = None,
        metadata: Optional["MetadataLike"] = None
    ) -> "SignResponse":
        return await self._unary_unary(
            "/keyring.v1.KeyringService/Sign",
            sign_request,
            SignResponse,
            timeout=timeout,
            deadline=deadline,
            metadata=metadata,
        )

    async def save_offline(
        self,
        save_offline_request: "SaveOfflineRequest",
        *,
        timeout: Optional[float] = None,
        deadline: Optional["Deadline"] = None,
        metadata: Optional["MetadataLike"] = None
    ) -> "SaveOfflineResponse":
        return await self._unary_unary(
            "/keyring.v1.KeyringService/SaveOffline",
            save_offline_request,
            SaveOfflineResponse,
            timeout=timeout,
            deadline=deadline,
            metadata=metadata,
        )


class KeyringServiceBase(ServiceBase):
    async def backend(self, backend_request: "BackendRequest") -> "BackendResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def key(self, key_request: "KeyRequest") -> "KeyResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def new_account(
        self, new_account_request: "NewAccountRequest"
    ) -> "NewAccountResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def sign(self, sign_request: "SignRequest") -> "SignResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def save_offline(
        self, save_offline_request: "SaveOfflineRequest"
    ) -> "SaveOfflineResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def __rpc_backend(
        self, stream: "grpclib.server.Stream[BackendRequest, BackendResponse]"
    ) -> None:
        request = await stream.recv_message()
        response = await self.backend(request)
        await stream.send_message(response)

    async def __rpc_key(
        self, stream: "grpclib.server.Stream[KeyRequest, KeyResponse]"
    ) -> None:
        request = await stream.recv_message()
        response = await self.key(request)
        await stream.send_message(response)

    async def __rpc_new_account(
        self, stream: "grpclib.server.Stream[NewAccountRequest, NewAccountResponse]"
    ) -> None:
        request = await stream.recv_message()
        response = await self.new_account(request)
        await stream.send_message(response)

    async def __rpc_sign(
        self, stream: "grpclib.server.Stream[SignRequest, SignResponse]"
    ) -> None:
        request = await stream.recv_message()
        response = await self.sign(request)
        await stream.send_message(response)

    async def __rpc_save_offline(
        self, stream: "grpclib.server.Stream[SaveOfflineRequest, SaveOfflineResponse]"
    ) -> None:
        request = await stream.recv_message()
        response = await self.save_offline(request)
        await stream.send_message(response)

    def __mapping__(self) -> Dict[str, grpclib.const.Handler]:
        return {
            "/keyring.v1.KeyringService/Backend": grpclib.const.Handler(
                self.__rpc_backend,
                grpclib.const.Cardinality.UNARY_UNARY,
                BackendRequest,
                BackendResponse,
            ),
            "/keyring.v1.KeyringService/Key": grpclib.const.Handler(
                self.__rpc_key,
                grpclib.const.Cardinality.UNARY_UNARY,
                KeyRequest,
                KeyResponse,
            ),
            "/keyring.v1.KeyringService/NewAccount": grpclib.const.Handler(
                self.__rpc_new_account,
                grpclib.const.Cardinality.UNARY_UNARY,
                NewAccountRequest,
                NewAccountResponse,
            ),
            "/keyring.v1.KeyringService/Sign": grpclib.const.Handler(
                self.__rpc_sign,
                grpclib.const.Cardinality.UNARY_UNARY,
                SignRequest,
                SignResponse,
            ),
            "/keyring.v1.KeyringService/SaveOffline": grpclib.const.Handler(
                self.__rpc_save_offline,
                grpclib.const.Cardinality.UNARY_UNARY,
                SaveOfflineRequest,
                SaveOfflineResponse,
            ),
        }
