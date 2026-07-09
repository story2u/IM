import base64
import hashlib
import socket
import struct
from typing import Any

import httpx
import xmltodict
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from redis.asyncio import Redis

from app.core.config import Settings
from app.domain.enums import IMChannel
from app.domain.ports import InboundMessage, SendReceipt


class WeComCryptoError(ValueError):
    pass


class WeComCrypto:
    def __init__(self, token: str, encoding_aes_key: str, receive_id: str) -> None:
        self.token = token
        self.receive_id = receive_id
        self.key = base64.b64decode(f"{encoding_aes_key}=")
        if len(self.key) != 32:
            raise WeComCryptoError("invalid wecom aes key length")
        self.iv = self.key[:16]

    def verify_signature(self, signature: str, timestamp: str, nonce: str, encrypted: str) -> None:
        values = [self.token, timestamp, nonce, encrypted]
        digest = hashlib.sha1("".join(sorted(values)).encode("utf-8")).hexdigest()
        if digest != signature:
            raise WeComCryptoError("invalid wecom signature")

    def decrypt(self, encrypted: str) -> str:
        encrypted_bytes = base64.b64decode(encrypted)
        cipher = Cipher(algorithms.AES(self.key), modes.CBC(self.iv))
        decryptor = cipher.decryptor()
        padded = decryptor.update(encrypted_bytes) + decryptor.finalize()
        content = self._pkcs7_unpad(padded)
        msg_len = socket.ntohl(struct.unpack("I", content[16:20])[0])
        message = content[20 : 20 + msg_len]
        receive_id = content[20 + msg_len :].decode("utf-8")
        if self.receive_id and receive_id not in {self.receive_id, self.receive_id.lower()}:
            raise WeComCryptoError("invalid wecom receive id")
        return message.decode("utf-8")

    def _pkcs7_unpad(self, value: bytes) -> bytes:
        pad = value[-1]
        if pad < 1 or pad > 32:
            raise WeComCryptoError("invalid pkcs7 padding")
        return value[:-pad]


class WeComAdapter:
    channel = IMChannel.WECOM

    def __init__(self, settings: Settings, redis: Redis | None = None) -> None:
        self.settings = settings
        self.redis = redis

    async def verify_url(self, query: dict[str, str]) -> str:
        encrypted = query.get("echostr", "")
        crypto = self._crypto()
        crypto.verify_signature(
            signature=query.get("msg_signature", ""),
            timestamp=query.get("timestamp", ""),
            nonce=query.get("nonce", ""),
            encrypted=encrypted,
        )
        return crypto.decrypt(encrypted)

    async def parse_webhook(
        self,
        payload: dict[str, Any],
        headers: dict[str, str],
        query: dict[str, str] | None = None,
    ) -> InboundMessage | None:
        query = query or {}
        encrypted = self._extract_encrypt(payload)
        crypto = self._crypto()
        crypto.verify_signature(
            signature=query.get("msg_signature", ""),
            timestamp=query.get("timestamp", ""),
            nonce=query.get("nonce", ""),
            encrypted=encrypted,
        )
        decrypted_xml = crypto.decrypt(encrypted)
        body = xmltodict.parse(decrypted_xml).get("xml", {})

        msg_type = body.get("MsgType")
        if msg_type != "text":
            return None

        content = body.get("Content")
        if not content:
            return None

        external_message_id = str(body.get("MsgId") or f"{body.get('FromUserName')}-{body.get('CreateTime')}")
        conversation_id = str(body.get("FromUserName"))

        return InboundMessage(
            channel=self.channel,
            external_message_id=external_message_id,
            conversation_id=conversation_id,
            sender_external_id=str(body.get("FromUserName")),
            sender_display_name=str(body.get("FromUserName")),
            text=str(content),
            raw_payload={"encrypted": payload, "decrypted": body},
        )

    async def send_message(self, conversation_id: str, text: str) -> SendReceipt:
        if not self.settings.im_send_enabled:
            return SendReceipt(
                provider_message_id=None,
                raw_response={"dry_run": True, "channel": self.channel, "touser": conversation_id},
            )

        token = await self._access_token()
        async with httpx.AsyncClient(timeout=10.0) as client:
            response = await client.post(
                f"https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={token}",
                json={
                    "touser": conversation_id,
                    "msgtype": "text",
                    "agentid": int(self.settings.wecom_agent_id),
                    "text": {"content": text},
                    "safe": 0,
                },
            )
            response.raise_for_status()
            data = response.json()
            if data.get("errcode") != 0:
                raise RuntimeError(f"wecom send failed: {data}")
            return SendReceipt(provider_message_id=None, raw_response=data)

    async def _access_token(self) -> str:
        cache_key = "wecom:access_token"
        if self.redis:
            cached = await self.redis.get(cache_key)
            if cached:
                return str(cached)

        async with httpx.AsyncClient(timeout=10.0) as client:
            response = await client.get(
                "https://qyapi.weixin.qq.com/cgi-bin/gettoken",
                params={"corpid": self.settings.wecom_corp_id, "corpsecret": self.settings.wecom_secret},
            )
            response.raise_for_status()
            data = response.json()
            if data.get("errcode") != 0:
                raise RuntimeError(f"wecom token failed: {data}")
            token = data["access_token"]

        if self.redis:
            await self.redis.set(cache_key, token, ex=7000)
        return token

    def _crypto(self) -> WeComCrypto:
        return WeComCrypto(
            token=self.settings.wecom_token,
            encoding_aes_key=self.settings.wecom_aes_key,
            receive_id=self.settings.wecom_corp_id,
        )

    def _extract_encrypt(self, payload: dict[str, Any]) -> str:
        if "xml" in payload:
            xml_node = payload["xml"]
            if isinstance(xml_node, dict) and xml_node.get("Encrypt"):
                return str(xml_node["Encrypt"])
        if payload.get("Encrypt"):
            return str(payload["Encrypt"])
        raise WeComCryptoError("missing Encrypt field")
