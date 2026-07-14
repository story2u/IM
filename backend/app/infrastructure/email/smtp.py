import asyncio
import html
import smtplib
import ssl
from email.message import EmailMessage
from urllib.parse import quote

from app.core.config import Settings


class SMTPEmailSender:
    def __init__(self, settings: Settings) -> None:
        if not settings.password_reset_email_configured:
            raise RuntimeError("password reset email is not configured")
        self.settings = settings

    async def send_password_reset(self, *, recipient: str, token: str, code: str) -> None:
        await asyncio.to_thread(
            self._send_password_reset,
            recipient=recipient,
            token=token,
            code=code,
        )

    def _send_password_reset(self, *, recipient: str, token: str, code: str) -> None:
        reset_url = (
            f"{self.settings.frontend_base_url.rstrip('/')}/reset-password"
            f"?token={quote(token, safe='')}"
        )
        message = EmailMessage()
        message["Subject"] = "重置商机雷达密码"
        message["From"] = (
            f"{self.settings.smtp_from_name} <{self.settings.smtp_from_email}>"
        )
        message["To"] = recipient
        message.set_content(
            "你正在重置商机雷达密码。\n\n"
            f"验证码：{code}\n"
            f"重置链接：{reset_url}\n\n"
            f"以上凭据将在 {self.settings.password_reset_ttl_minutes} 分钟后失效。"
            "如果不是你本人操作，请忽略此邮件。"
        )
        message.add_alternative(
            "<p>你正在重置商机雷达密码。</p>"
            f"<p>验证码：<strong>{html.escape(code)}</strong></p>"
            f'<p><a href="{html.escape(reset_url, quote=True)}">重置密码</a></p>'
            f"<p>以上凭据将在 {self.settings.password_reset_ttl_minutes} 分钟后失效。"
            "如果不是你本人操作，请忽略此邮件。</p>",
            subtype="html",
        )

        tls_context = ssl.create_default_context()
        if self.settings.smtp_use_tls:
            client_context = smtplib.SMTP_SSL(
                self.settings.smtp_host,
                self.settings.smtp_port,
                timeout=self.settings.smtp_timeout_seconds,
                context=tls_context,
            )
        else:
            client_context = smtplib.SMTP(
                self.settings.smtp_host,
                self.settings.smtp_port,
                timeout=self.settings.smtp_timeout_seconds,
            )
        with client_context as client:
            if self.settings.smtp_starttls:
                client.starttls(context=tls_context)
            if self.settings.smtp_username:
                client.login(self.settings.smtp_username, self.settings.smtp_password)
            client.send_message(message)
