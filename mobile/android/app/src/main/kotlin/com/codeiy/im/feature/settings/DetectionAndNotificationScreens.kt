package com.codeiy.im.feature.settings

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Close
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.codeiy.im.model.DetectionSettings
import com.codeiy.im.model.NotificationSettings
import com.codeiy.im.ui.theme.AppColors

/** 商机识别规则：关键词增删 + AI 语义开关，保存失败回滚。 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DetectionSettingsScreen(model: SettingsViewModel, detection: DetectionSettings, onBack: () -> Unit) {
    val keywords = remember { mutableStateListOf<String>().apply { addAll(detection.keywords) } }
    var aiEnabled by remember { mutableStateOf(detection.aiSemanticsEnabled) }
    var newKeyword by remember { mutableStateOf("") }
    var error by remember { mutableStateOf<String?>(null) }

    fun addKeyword() {
        val trimmed = newKeyword.trim()
        if (trimmed.isNotEmpty() && trimmed !in keywords) keywords.add(trimmed)
        newKeyword = ""
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("商机识别规则") },
                navigationIcon = { IconButton(onClick = onBack) { Icon(Icons.AutoMirrored.Filled.ArrowBack, "返回") } },
                actions = {
                    TextButton(onClick = {
                        error = null
                        model.saveDetection(keywords.toList(), aiEnabled) { error = it }
                    }) { Text("保存") }
                },
            )
        },
    ) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Column(Modifier.weight(1f)) {
                    Text("AI 语义识别", style = MaterialTheme.typography.bodyLarge)
                    Text("开启后除关键词外，AI 会理解语义识别潜在商机。", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                }
                Switch(checked = aiEnabled, onCheckedChange = { aiEnabled = it })
            }

            Text("关键词", style = MaterialTheme.typography.titleSmall)
            if (keywords.isEmpty()) {
                Text("暂无关键词", color = AppColors.muted)
            } else {
                keywords.forEach { keyword ->
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(keyword, Modifier.weight(1f))
                        IconButton(onClick = { keywords.remove(keyword) }) {
                            Icon(Icons.Filled.Close, contentDescription = "删除 $keyword", tint = AppColors.muted)
                        }
                    }
                }
            }
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp), verticalAlignment = Alignment.CenterVertically) {
                OutlinedTextField(value = newKeyword, onValueChange = { newKeyword = it }, placeholder = { Text("添加关键词") }, modifier = Modifier.weight(1f), singleLine = true)
                TextButton(enabled = newKeyword.trim().isNotEmpty(), onClick = ::addKeyword) { Text("添加") }
            }
            error?.let { Text(it, color = AppColors.destructive, style = MaterialTheme.typography.labelSmall) }
        }
    }
}

/** 通知偏好：4 开关，推送未开放时标注"启用后生效"，开关即时保存失败回滚。 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NotificationSettingsScreen(
    model: SettingsViewModel,
    notifications: NotificationSettings,
    pushAvailable: Boolean,
    onBack: () -> Unit,
) {
    var prefs by remember { mutableStateOf(notifications) }
    var error by remember { mutableStateOf<String?>(null) }

    fun update(next: NotificationSettings) {
        prefs = next
        error = null
        model.saveNotifications(next) { message -> error = message; prefs = notifications }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("通知设置") },
                navigationIcon = { IconButton(onClick = onBack) { Icon(Icons.AutoMirrored.Filled.ArrowBack, "返回") } },
            )
        },
    ) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
            if (!pushAvailable) {
                Text("推送服务尚未开放，偏好会保存，将在启用后生效。", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
            }
            toggleRow("新商机提醒", prefs.newOpportunityEnabled) { update(prefs.copy(newOpportunityEnabled = it)) }
            toggleRow("AI 已回复通知", prefs.aiRepliedEnabled) { update(prefs.copy(aiRepliedEnabled = it)) }
            toggleRow("每日商机摘要", prefs.dailyDigestEnabled) { update(prefs.copy(dailyDigestEnabled = it)) }
            toggleRow("仅紧急商机", prefs.urgentOnly) { update(prefs.copy(urgentOnly = it)) }
            error?.let { Text(it, color = AppColors.destructive, style = MaterialTheme.typography.labelSmall) }
        }
    }
}

@Composable
private fun toggleRow(title: String, checked: Boolean, onChange: (Boolean) -> Unit) {
    Row(Modifier.fillMaxWidth(), verticalAlignment = Alignment.CenterVertically) {
        Text(title, Modifier.weight(1f))
        Switch(checked = checked, onCheckedChange = onChange)
    }
}
