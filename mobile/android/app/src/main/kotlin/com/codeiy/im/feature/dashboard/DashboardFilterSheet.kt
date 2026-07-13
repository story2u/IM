package com.codeiy.im.feature.dashboard

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Check
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilterChip
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.codeiy.im.ui.theme.AppColors
import com.codeiy.im.ui.theme.SopStage
import com.codeiy.im.ui.theme.TrustLevel
import com.codeiy.im.ui.theme.color

/** 高级筛选 Bottom Sheet：草稿模式，点"应用"才写回。对齐 Web filter-panel.tsx。 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DashboardFilterSheet(
    query: DashboardQuery,
    keywordOptions: List<String>,
    onApply: (DashboardQuery) -> Unit,
    onDismiss: () -> Unit,
) {
    var draft by remember { mutableStateOf(query) }

    ModalBottomSheet(onDismissRequest = onDismiss) {
        Column(Modifier.fillMaxWidth().verticalScroll(rememberScrollState()).padding(horizontal = 16.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text("高级筛选", style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.weight(1f))
                if (draft.activeAdvancedCount > 0) {
                    Text("${draft.activeAdvancedCount} 项", style = MaterialTheme.typography.labelMedium, color = AppColors.muted)
                }
            }
            Spacer(Modifier.size(8.dp))

            section("时间范围") {
                chipRow(DashboardQuery.TimeRange.entries.map { it to it.label }, selected = draft.timeRange) {
                    draft = draft.copy(timeRange = it)
                }
                if (draft.timeRange == DashboardQuery.TimeRange.CUSTOM && !draft.customRangeValid) {
                    Text("开始日期不能晚于结束日期", style = MaterialTheme.typography.labelSmall, color = AppColors.destructive)
                }
            }

            section("消息来源") {
                chipRow(DashboardQuery.Source.entries.map { it to it.label }, selected = draft.source) {
                    draft = draft.copy(source = it)
                }
            }

            section("可信度") {
                multiChipRow(TrustLevel.entries.map { it to it.label }, selected = draft.trustLevels) { level, on ->
                    draft = draft.copy(trustLevels = draft.trustLevels.toggle(level, on))
                }
            }

            section("流程阶段") {
                multiChipRow(SopStage.entries.map { it to it.label }, selected = draft.stages) { stage, on ->
                    draft = draft.copy(stages = draft.stages.toggle(stage, on))
                }
            }

            if (keywordOptions.isNotEmpty()) {
                section("关键词标签") {
                    multiChipRow(keywordOptions.map { it to it }, selected = draft.keywords) { keyword, on ->
                        draft = draft.copy(keywords = draft.keywords.toggle(keyword, on))
                    }
                }
            }

            HorizontalDivider(Modifier.padding(vertical = 8.dp))
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                TextButton(onClick = onDismiss) { Text("取消") }
                TextButton(onClick = { draft = resetAdvanced(draft) }) { Text("重置") }
                Spacer(Modifier.weight(1f))
                TextButton(enabled = draft.customRangeValid, onClick = { onApply(draft) }) { Text("应用") }
            }
            Spacer(Modifier.size(24.dp))
        }
    }
}

@Composable
private fun section(title: String, content: @Composable () -> Unit) {
    Text(title, style = MaterialTheme.typography.titleSmall, modifier = Modifier.padding(top = 12.dp, bottom = 4.dp))
    content()
}

@OptIn(ExperimentalMaterial3Api::class, androidx.compose.foundation.layout.ExperimentalLayoutApi::class)
@Composable
private fun <T> chipRow(options: List<Pair<T, String>>, selected: T, onSelect: (T) -> Unit) {
    androidx.compose.foundation.layout.FlowRow(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
        options.forEach { (value, label) ->
            FilterChip(selected = value == selected, onClick = { onSelect(value) }, label = { Text(label) })
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class, androidx.compose.foundation.layout.ExperimentalLayoutApi::class)
@Composable
private fun <T> multiChipRow(options: List<Pair<T, String>>, selected: Set<T>, onToggle: (T, Boolean) -> Unit) {
    androidx.compose.foundation.layout.FlowRow(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
        options.forEach { (value, label) ->
            val isOn = value in selected
            FilterChip(
                selected = isOn,
                onClick = { onToggle(value, !isOn) },
                label = { Text(label) },
                leadingIcon = if (isOn) {
                    { Icon(Icons.Filled.Check, contentDescription = null, modifier = Modifier.size(16.dp)) }
                } else null,
            )
        }
    }
}

private fun <T> Set<T>.toggle(item: T, on: Boolean): Set<T> = if (on) this + item else this - item

private fun resetAdvanced(q: DashboardQuery) = q.copy(
    source = DashboardQuery.Source.ALL,
    timeRange = DashboardQuery.TimeRange.ALL,
    customFrom = null,
    customTo = null,
    trustLevels = emptySet(),
    stages = emptySet(),
    keywords = emptySet(),
)
