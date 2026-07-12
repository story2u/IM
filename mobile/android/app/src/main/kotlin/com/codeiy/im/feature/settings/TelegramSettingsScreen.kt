package com.codeiy.im.feature.settings

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Group
import androidx.compose.material.icons.filled.Person
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import com.codeiy.im.core.auth.SessionStore
import com.codeiy.im.core.network.RadarApi
import com.codeiy.im.core.network.api
import com.codeiy.im.model.TelegramConnectionDTO
import com.codeiy.im.model.TelegramConnectionEnabledUpdate
import com.codeiy.im.model.TelegramConnectionHealth
import com.codeiy.im.ui.theme.AppColors
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class TelegramSettingsViewModel(private val service: RadarApi) : ViewModel() {
    data class UiState(
        val health: TelegramConnectionHealth? = null,
        val connections: List<TelegramConnectionDTO> = emptyList(),
        val isLoading: Boolean = false,
        val error: String? = null,
    )

    private val _state = MutableStateFlow(UiState())
    val state: StateFlow<UiState> = _state

    fun load() {
        _state.value = _state.value.copy(isLoading = true)
        viewModelScope.launch {
            try {
                coroutineScope {
                    val health = async { api { service.telegramHealth() } }
                    val connections = async { api { service.telegramConnections() } }
                    _state.value = UiState(health = health.await(), connections = connections.await(), isLoading = false)
                }
            } catch (e: Exception) {
                _state.value = _state.value.copy(isLoading = false, error = e.message)
            }
        }
    }

    fun toggle(connection: TelegramConnectionDTO) {
        viewModelScope.launch {
            try {
                val updated = api { service.setTelegramConnectionEnabled(connection.id, TelegramConnectionEnabledUpdate(!connection.enabled)) }
                _state.value = _state.value.copy(connections = _state.value.connections.map { if (it.id == updated.id) updated else it })
            } catch (e: Exception) {
                _state.value = _state.value.copy(error = e.message)
            }
        }
    }
}

/** Telegram 连接：真实读取 health + connections + sources，可停用/启用。
 *  连接建立向导（Bot/Business/QR 握手+深链+轮询）为后续迭代，未配置能力如实标注。 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TelegramSettingsScreen(session: SessionStore, onBack: () -> Unit) {
    val model: TelegramSettingsViewModel = viewModel { TelegramSettingsViewModel(session.api.service) }
    val state by model.state.collectAsStateWithLifecycle()

    LaunchedEffect(Unit) { model.load() }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Telegram 连接") },
                navigationIcon = { IconButton(onClick = onBack) { Icon(Icons.AutoMirrored.Filled.ArrowBack, "返回") } },
            )
        },
    ) { padding ->
        LazyColumn(Modifier.padding(padding).fillMaxSize(), contentPadding = androidx.compose.foundation.layout.PaddingValues(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            state.health?.let { health ->
                item {
                    Column(verticalArrangement = Arrangement.spacedBy(4.dp)) {
                        Text("服务状态", style = MaterialTheme.typography.titleSmall)
                        statusRow("Bot", if (health.botConfigured) (health.botUsername?.let { "@$it" } ?: "已配置") else "管理员尚未配置")
                        statusRow("Business 私聊", if (health.businessAvailable) "可用" else "管理员尚未配置")
                        statusRow("普通账号 QR", if (health.mtprotoQrAvailable) "可用" else "管理员尚未配置")
                    }
                }
            }

            if (state.connections.isEmpty() && !state.isLoading) {
                item {
                    Column(Modifier.fillMaxWidth().padding(16.dp), horizontalAlignment = Alignment.CenterHorizontally) {
                        Text("尚无连接", style = MaterialTheme.typography.bodyLarge)
                        Text("连接建立向导即将上线；当前可在 Web 端完成绑定后在此管理。", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                    }
                }
            }

            items(state.connections.size) { index ->
                val connection = state.connections[index]
                Column {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(connection.label, Modifier.weight(1f), style = MaterialTheme.typography.titleSmall)
                        Switch(checked = connection.enabled, onCheckedChange = { model.toggle(connection) })
                    }
                    Text("状态：${connection.status}", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                    connection.lastError?.let { Text(it, style = MaterialTheme.typography.labelSmall, color = AppColors.warning) }
                    connection.sources.forEach { source ->
                        Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.padding(vertical = 4.dp)) {
                            Icon(if (source.sourceType == "private") Icons.Filled.Person else Icons.Filled.Group, contentDescription = null, tint = AppColors.muted, modifier = Modifier.size(16.dp))
                            Column(Modifier.weight(1f)) {
                                Text(source.displayName)
                                if (source.quotaPaused && source.quotaReason != null) {
                                    Text(source.quotaReason!!, style = MaterialTheme.typography.labelSmall, color = AppColors.warning)
                                }
                            }
                            if (!source.enabled) Text("已停用", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                        }
                    }
                    HorizontalDivider(Modifier.padding(top = 8.dp))
                }
            }

            state.error?.let { item { Text(it, color = AppColors.destructive, style = MaterialTheme.typography.labelSmall) } }
        }
    }
}

@Composable
private fun statusRow(label: String, value: String) {
    Row {
        Text(label, Modifier.weight(1f), color = AppColors.muted, style = MaterialTheme.typography.bodyMedium)
        Text(value, style = MaterialTheme.typography.bodyMedium)
    }
}
