package com.codeiy.im.feature.dashboard

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.GridItemSpan
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.clickable
import androidx.compose.foundation.background
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material.icons.filled.Tune
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilterChip
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SegmentedButton
import androidx.compose.material3.SegmentedButtonDefaults
import androidx.compose.material3.SingleChoiceSegmentedButtonRow
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.pulltorefresh.PullToRefreshBox
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import com.codeiy.im.core.auth.SessionStore
import com.codeiy.im.model.FrontendOpportunityStatus
import com.codeiy.im.model.IMChannel
import com.codeiy.im.ui.theme.AppCard
import com.codeiy.im.ui.theme.AppColors

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DashboardScreen(session: SessionStore, onOpenOpportunity: (String) -> Unit) {
    val model: DashboardViewModel = viewModel { DashboardViewModel(session.api.service) }
    val state by model.state.collectAsStateWithLifecycle()
    var showFilterSheet by remember { mutableStateOf(false) }
    var showPlatformMenu by remember { mutableStateOf(false) }
    var showSortMenu by remember { mutableStateOf(false) }

    LaunchedEffect(Unit) {
        // 用用户时区计算"今天"边界（失败或未设置时用设备时区）。
        runCatching { com.codeiy.im.core.network.api { session.api.service.settings() } }
            .getOrNull()?.let { model.setTimezone(it.workSchedule.timezone) }
        model.refresh()
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Column {
                        Text("商机看板", style = MaterialTheme.typography.titleLarge)
                        Text("${state.pendingCount} 条待处理", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                    }
                },
                actions = {
                    IconButton(onClick = { model.refresh() }) { Icon(Icons.Filled.Refresh, contentDescription = "刷新") }
                },
            )
        },
    ) { padding ->
        PullToRefreshBox(
            isRefreshing = state.isInitialLoading && state.items.isNotEmpty(),
            onRefresh = { model.refresh() },
            modifier = Modifier.padding(padding).fillMaxSize(),
        ) {
            LazyVerticalGrid(
                columns = GridCells.Adaptive(320.dp),
                modifier = Modifier.fillMaxSize(),
                contentPadding = androidx.compose.foundation.layout.PaddingValues(12.dp),
                verticalArrangement = Arrangement.spacedBy(12.dp),
                horizontalArrangement = Arrangement.spacedBy(12.dp),
            ) {
                item(span = { GridItemSpan(maxLineSpan) }) {
                    Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                        if (state.attentionItems.isNotEmpty()) attentionBanner(state, onOpenOpportunity)
                        primaryFilters(
                            state = state,
                            model = model,
                            showPlatformMenu = showPlatformMenu,
                            onPlatformMenu = { showPlatformMenu = it },
                            showSortMenu = showSortMenu,
                            onSortMenu = { showSortMenu = it },
                            onAdvanced = { showFilterSheet = true },
                        )
                        resultSummary(state, model)
                    }
                }
                dashboardBody(state, model, onOpenOpportunity)
            }
        }
    }

    if (showFilterSheet) {
        DashboardFilterSheet(
            query = state.query,
            keywordOptions = state.keywordOptions,
            onApply = { model.setQuery(it); showFilterSheet = false },
            onDismiss = { showFilterSheet = false },
        )
    }
}

private fun androidx.compose.foundation.lazy.grid.LazyGridScope.dashboardBody(
    state: DashboardViewModel.UiState,
    model: DashboardViewModel,
    onOpen: (String) -> Unit,
) {
    when {
        state.initialError != null && state.items.isEmpty() -> item(span = { GridItemSpan(maxLineSpan) }) {
            Column(Modifier.fillMaxWidth().padding(24.dp), horizontalAlignment = Alignment.CenterHorizontally) {
                Text(state.initialError, color = AppColors.destructive)
                TextButton(onClick = { model.retryInitial() }) { Text("重试") }
            }
        }
        state.items.isEmpty() && state.isInitialLoading -> item(span = { GridItemSpan(maxLineSpan) }) {
            Box(Modifier.fillMaxWidth().padding(32.dp), contentAlignment = Alignment.Center) { CircularProgressIndicator() }
        }
        state.items.isEmpty() -> item(span = { GridItemSpan(maxLineSpan) }) {
            Box(Modifier.fillMaxWidth().padding(32.dp), contentAlignment = Alignment.Center) {
                Text("暂无匹配的商机", color = AppColors.muted)
            }
        }
        else -> {
            if (state.pageError != null) {
                item(span = { GridItemSpan(maxLineSpan) }) {
                    Text("刷新失败：${state.pageError}（点按重试）", color = AppColors.destructive, style = MaterialTheme.typography.labelSmall, modifier = Modifier.fillMaxWidth().clickable { model.refresh() }.padding(6.dp))
                }
            }
            items(state.items, key = { it.id }) { opportunity ->
                OpportunityCard(opportunity, onClick = { onOpen(opportunity.id) })
                LaunchedEffect(opportunity.id) { model.loadMoreIfNeeded(opportunity) }
            }
            if (state.isLoadingMore) {
                item(span = { GridItemSpan(maxLineSpan) }) {
                    Box(Modifier.fillMaxWidth().padding(16.dp), contentAlignment = Alignment.Center) { CircularProgressIndicator() }
                }
            }
        }
    }
}

@Composable
private fun attentionBanner(state: DashboardViewModel.UiState, onOpen: (String) -> Unit) {
    AppCard {
        Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                Icon(Icons.Filled.Warning, contentDescription = null, tint = AppColors.warning, modifier = Modifier.size(18.dp))
                Text("pi Agent 发现 ${state.attentionItems.size} 条重大商机", style = MaterialTheme.typography.titleSmall, color = AppColors.warning)
            }
            Text("请优先核对链接结论和后续行动建议，外部动作仍需人工批准。", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
            state.attentionItems.take(3).forEach { item ->
                Text(item.contactName, style = MaterialTheme.typography.bodyMedium, modifier = Modifier.fillMaxWidth().clickable { onOpen(item.id) })
            }
            if (state.attentionItems.size > 3) {
                Text("查看全部", style = MaterialTheme.typography.labelSmall, color = MaterialTheme.colorScheme.primary)
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun primaryFilters(
    state: DashboardViewModel.UiState,
    model: DashboardViewModel,
    showPlatformMenu: Boolean,
    onPlatformMenu: (Boolean) -> Unit,
    showSortMenu: Boolean,
    onSortMenu: (Boolean) -> Unit,
    onAdvanced: () -> Unit,
) {
    val q = state.query
    Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
        val statuses = listOf<Pair<FrontendOpportunityStatus?, String>>(
            null to "全部",
            FrontendOpportunityStatus.PENDING to "待处理",
            FrontendOpportunityStatus.REPLIED to "已回复",
            FrontendOpportunityStatus.IGNORED to "已忽略",
        )
        SingleChoiceSegmentedButtonRow(Modifier.fillMaxWidth()) {
            statuses.forEachIndexed { index, (value, label) ->
                SegmentedButton(
                    selected = q.status == value,
                    onClick = { model.setQuery(q.copy(status = value)) },
                    shape = SegmentedButtonDefaults.itemShape(index, statuses.size),
                ) { Text(label) }
            }
        }
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp), verticalAlignment = Alignment.CenterVertically) {
            Box {
                FilterChip(selected = q.platform != null, onClick = { onPlatformMenu(true) }, label = { Text(q.platform?.label ?: "全部平台") })
                DropdownMenu(expanded = showPlatformMenu, onDismissRequest = { onPlatformMenu(false) }) {
                    DropdownMenuItem(text = { Text("全部平台") }, onClick = { model.setQuery(q.copy(platform = null)); onPlatformMenu(false) })
                    listOf(IMChannel.TELEGRAM, IMChannel.WECOM).forEach { channel ->
                        DropdownMenuItem(text = { Text(channel.label) }, onClick = { model.setQuery(q.copy(platform = channel)); onPlatformMenu(false) })
                    }
                }
            }
            Box {
                FilterChip(selected = true, onClick = { onSortMenu(true) }, label = { Text(q.sort.label) })
                DropdownMenu(expanded = showSortMenu, onDismissRequest = { onSortMenu(false) }) {
                    DashboardQuery.Sort.entries.forEach { sort ->
                        DropdownMenuItem(text = { Text(sort.label) }, onClick = { model.setQuery(q.copy(sort = sort)); onSortMenu(false) })
                    }
                }
            }
            Spacer(Modifier.weight(1f))
            IconButton(onClick = onAdvanced) {
                Box(contentAlignment = Alignment.TopEnd) {
                    Icon(Icons.Filled.Tune, contentDescription = "高级筛选")
                    if (q.activeAdvancedCount > 0) {
                        Box(Modifier.size(14.dp).background(MaterialTheme.colorScheme.primary, CircleShape), contentAlignment = Alignment.Center) {
                            Text("${q.activeAdvancedCount}", style = MaterialTheme.typography.labelSmall, color = MaterialTheme.colorScheme.onPrimary)
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun resultSummary(state: DashboardViewModel.UiState, model: DashboardViewModel) {
    val q = state.query
    Row(verticalAlignment = Alignment.CenterVertically) {
        Text("共 ${state.total} 条商机", style = MaterialTheme.typography.labelMedium, color = AppColors.muted)
        Spacer(Modifier.weight(1f))
        if (q.activeAdvancedCount > 0 || q.status != null || q.platform != null) {
            TextButton(onClick = { model.setQuery(DashboardQuery(timezoneId = q.timezoneId)) }) { Text("清空筛选") }
        }
    }
}
