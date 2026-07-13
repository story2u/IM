package com.codeiy.im.feature.dashboard

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.codeiy.im.core.network.RadarApi
import com.codeiy.im.core.network.api
import com.codeiy.im.model.DashboardResponse
import com.codeiy.im.model.Opportunity
import kotlin.coroutines.cancellation.CancellationException
import kotlinx.coroutines.Job
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class DashboardViewModel(private val service: RadarApi) : ViewModel() {
    data class UiState(
        val items: List<Opportunity> = emptyList(),
        val total: Int = 0,
        val pendingCount: Int = 0,
        val attentionItems: List<Opportunity> = emptyList(),
        val keywordOptions: List<String> = emptyList(),
        val query: DashboardQuery = DashboardQuery(),
        val isInitialLoading: Boolean = false,
        val isLoadingMore: Boolean = false,
        val initialError: String? = null,
        val pageError: String? = null,
        val canLoadMore: Boolean = false,
    )

    private val _state = MutableStateFlow(UiState())
    val state: StateFlow<UiState> = _state
    private val pageSize = 20

    // 单一在途加载：新筛选取消旧请求，慢响应不会覆盖新结果。
    private var loadJob: Job? = null

    fun setQuery(query: DashboardQuery) {
        if (query == _state.value.query) return
        _state.value = _state.value.copy(query = query)
        refresh()
    }

    fun setTimezone(id: String) {
        _state.value = _state.value.copy(query = _state.value.query.copy(timezoneId = id))
    }

    fun refresh(): Job {
        loadJob?.cancel()
        val job = viewModelScope.launch { load(reset = true) }
        loadJob = job
        return job
    }

    fun retryInitial() {
        _state.value = _state.value.copy(initialError = null)
        refresh()
    }

    fun loadMoreIfNeeded(item: Opportunity) {
        val s = _state.value
        if (s.canLoadMore && loadJob?.isActive != true && item.id == s.items.lastOrNull()?.id) {
            loadJob = viewModelScope.launch { load(reset = false) }
        }
    }

    private suspend fun load(reset: Boolean) {
        val requested = _state.value.query
        _state.value = _state.value.copy(
            isInitialLoading = reset && _state.value.items.isEmpty(),
            isLoadingMore = !reset,
        )
        try {
            val bounds = requested.resolvedBounds()
            val offset = if (reset) 0 else _state.value.items.size
            val response = fetch(requested, bounds, offset)
            // 写入前校验筛选未变（慢响应防覆盖）。
            if (_state.value.query != requested) return
            val merged = if (reset) response.items else _state.value.items + response.items
            _state.value = _state.value.copy(
                items = merged,
                total = response.total,
                pendingCount = response.pendingCount,
                attentionItems = response.attentionItems,
                keywordOptions = response.keywordOptions,
                canLoadMore = merged.size < response.total && response.items.isNotEmpty(),
                isInitialLoading = false,
                isLoadingMore = false,
                initialError = null,
                pageError = null,
            )
        } catch (e: CancellationException) {
            throw e
        } catch (e: Exception) {
            // 已有数据时失败保留旧数据，仅提示；首屏失败才空态。
            if (_state.value.items.isEmpty()) {
                _state.value = _state.value.copy(isInitialLoading = false, initialError = e.message)
            } else {
                _state.value = _state.value.copy(isLoadingMore = false, pageError = e.message)
            }
        }
    }

    private suspend fun fetch(q: DashboardQuery, bounds: Pair<String?, String?>, offset: Int): DashboardResponse = api {
        service.dashboard(
            status = q.status?.let { statusSerial(it) },
            platform = q.platform?.let { platformSerial(it) },
            sourceType = q.source.serial,
            createdFrom = bounds.first,
            createdTo = bounds.second,
            trustLevels = q.trustLevels.map { it.serialName }.ifEmpty { null },
            sopStages = q.stages.map { it.key }.ifEmpty { null },
            keywords = q.keywords.toList().ifEmpty { null },
            sort = q.sort.serial,
            limit = pageSize,
            offset = offset,
        )
    }

    private fun statusSerial(value: com.codeiy.im.model.FrontendOpportunityStatus) = when (value) {
        com.codeiy.im.model.FrontendOpportunityStatus.PENDING -> "pending"
        com.codeiy.im.model.FrontendOpportunityStatus.REPLIED -> "replied"
        com.codeiy.im.model.FrontendOpportunityStatus.IGNORED -> "ignored"
        com.codeiy.im.model.FrontendOpportunityStatus.UNKNOWN -> "pending"
    }

    private fun platformSerial(value: com.codeiy.im.model.IMChannel) = when (value) {
        com.codeiy.im.model.IMChannel.TELEGRAM -> "telegram"
        com.codeiy.im.model.IMChannel.WECOM -> "wecom"
        com.codeiy.im.model.IMChannel.UNKNOWN -> "telegram"
    }
}
