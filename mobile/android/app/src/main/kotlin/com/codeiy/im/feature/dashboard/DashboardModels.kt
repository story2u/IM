package com.codeiy.im.feature.dashboard

import com.codeiy.im.model.FrontendOpportunityStatus
import com.codeiy.im.model.IMChannel
import com.codeiy.im.ui.theme.SopStage
import com.codeiy.im.ui.theme.TrustLevel
import java.time.LocalDate
import java.time.ZoneId
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

/** 看板筛选与排序。语义对齐 Web dashboard-filters.ts；"今天"用用户时区算 UTC 边界（无 MOCK_NOW）。 */
data class DashboardQuery(
    val status: FrontendOpportunityStatus? = null,
    val platform: IMChannel? = null,
    val source: Source = Source.ALL,
    val timeRange: TimeRange = TimeRange.ALL,
    val customFrom: LocalDate? = null,
    val customTo: LocalDate? = null,
    val trustLevels: Set<TrustLevel> = emptySet(),
    val stages: Set<SopStage> = emptySet(),
    val keywords: Set<String> = emptySet(),
    val sort: Sort = Sort.NEWEST,
    val timezoneId: String = ZoneId.systemDefault().id,
) {
    enum class Sort(val serial: String, val label: String) {
        NEWEST("newest", "最新优先"),
        OLDEST("oldest", "最早优先"),
        CONFIDENCE("confidence", "按相关度"),
        TRUST("trust", "按可信度"),
    }

    enum class TimeRange(val label: String) {
        ALL("全部时间"), TODAY("今天"), THREE_DAYS("近 3 天"), SEVEN_DAYS("近 7 天"), CUSTOM("自定义")
    }

    enum class Source(val serial: String?, val label: String) {
        ALL(null, "全部来源"), GROUP("group", "群消息"), PRIVATE("private", "私聊消息")
    }

    val activeAdvancedCount: Int
        get() {
            var n = 0
            if (source != Source.ALL) n++
            if (timeRange != TimeRange.ALL) n++
            if (keywords.isNotEmpty()) n++
            if (trustLevels.isNotEmpty()) n++
            if (stages.isNotEmpty()) n++
            return n
        }

    val customRangeValid: Boolean
        get() = timeRange != TimeRange.CUSTOM || customFrom == null || customTo == null || !customFrom.isAfter(customTo)

    /** 解析时间范围为 UTC ISO8601 起止（供后端 created_from/created_to）。 */
    fun resolvedBounds(): Pair<String?, String?> {
        val zone = runCatching { ZoneId.of(timezoneId) }.getOrDefault(ZoneId.systemDefault())
        val now = ZonedDateTime.now(zone)
        val iso = DateTimeFormatter.ISO_INSTANT
        fun startOfDay(date: LocalDate) = date.atStartOfDay(zone)
        return when (timeRange) {
            TimeRange.ALL -> null to null
            TimeRange.TODAY -> iso.format(now.toLocalDate().atStartOfDay(zone).toInstant()) to null
            TimeRange.THREE_DAYS -> iso.format(now.minusDays(3).toInstant()) to null
            TimeRange.SEVEN_DAYS -> iso.format(now.minusDays(7).toInstant()) to null
            TimeRange.CUSTOM -> {
                val from = customFrom?.let { iso.format(startOfDay(it).toInstant()) }
                val to = customTo?.let { iso.format(startOfDay(it).plusDays(1).minusSeconds(1).toInstant()) }
                from to to
            }
        }
    }
}
