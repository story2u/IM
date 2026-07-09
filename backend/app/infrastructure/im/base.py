from app.domain.enums import IMChannel
from app.domain.ports import IMAdapter


class AdapterRegistry:
    def __init__(self, adapters: list[IMAdapter]) -> None:
        self._adapters = {adapter.channel: adapter for adapter in adapters}

    def get(self, channel: IMChannel) -> IMAdapter:
        adapter = self._adapters.get(channel)
        if not adapter:
            raise KeyError(f"adapter not registered for channel={channel}")
        return adapter
