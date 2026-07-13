package com.codeiy.im.ui.theme

import androidx.compose.ui.graphics.Color
import com.codeiy.im.model.FrontendOpportunityStatus
import com.codeiy.im.model.IMChannel
import com.codeiy.im.model.OpportunityStatus
import com.codeiy.im.model.Priority

/** 跨端语义色，与 iOS AppColors / Web theme 语义一致（不追求像素一致）。 */
object AppColors {
    val success = Color(0xFF16A34A)
    val warning = Color(0xFFF59E0B)
    val destructive = Color(0xFFDC2626)
    val telegram = Color(0xFF2563EB)
    val wecom = Color(0xFF16A34A)
    val ai = Color(0xFF7C3AED)
    val muted = Color(0xFF6B7280)
}

/** 可信度等级：边界与后端 dashboard trust_levels、Web sop.ts 一致。 */
enum class TrustLevel(val label: String) {
    TRUSTED("安全可信"),
    UNVERIFIED("待核验"),
    SUSPICIOUS("可疑"),
    RISKY("高风险");

    companion object {
        fun from(score: Int): TrustLevel = when {
            score >= 80 -> TRUSTED
            score >= 60 -> UNVERIFIED
            score >= 40 -> SUSPICIOUS
            else -> RISKY
        }
    }

    val serialName: String
        get() = when (this) {
            TRUSTED -> "trusted"
            UNVERIFIED -> "unverified"
            SUSPICIOUS -> "suspicious"
            RISKY -> "risky"
        }
}

fun TrustLevel.color(): Color = when (this) {
    TrustLevel.TRUSTED -> AppColors.success
    TrustLevel.UNVERIFIED -> AppColors.muted
    TrustLevel.SUSPICIOUS -> AppColors.warning
    TrustLevel.RISKY -> AppColors.destructive
}

fun Priority.color(): Color = when (this) {
    Priority.URGENT -> AppColors.destructive
    Priority.HIGH -> AppColors.warning
    Priority.NORMAL -> Color(0xFF2563EB)
    Priority.LOW, Priority.UNKNOWN -> AppColors.muted
}

fun FrontendOpportunityStatus.color(): Color = when (this) {
    FrontendOpportunityStatus.PENDING -> AppColors.warning
    FrontendOpportunityStatus.REPLIED -> AppColors.success
    FrontendOpportunityStatus.IGNORED, FrontendOpportunityStatus.UNKNOWN -> AppColors.muted
}

fun OpportunityStatus.frontendColor(): Color = when (this) {
    OpportunityStatus.PENDING_HUMAN, OpportunityStatus.AI_AUTO_REPLY -> AppColors.warning
    OpportunityStatus.REPLIED, OpportunityStatus.FOLLOWING -> AppColors.success
    OpportunityStatus.IGNORED, OpportunityStatus.CLOSED, OpportunityStatus.UNKNOWN -> AppColors.muted
}

fun IMChannel.color(): Color = when (this) {
    IMChannel.TELEGRAM -> AppColors.telegram
    IMChannel.WECOM -> AppColors.wecom
    IMChannel.UNKNOWN -> AppColors.muted
}

/** SOP 流程阶段：键/标签/点色，对齐 Web sop.ts 与 iOS SopStage。 */
enum class SopStage(val key: String, val label: String, val dot: Color) {
    DETECTED("detected", "已发现", AppColors.muted),
    ANALYZING("analyzing", "AI 分析中", Color(0xFF2563EB)),
    VERIFIED("verified", "已核验", Color(0xFF2563EB)),
    CONTACT_EXTRACTED("contact_extracted", "已提取联系方式", Color(0xFF2563EB)),
    FRIEND_REQUESTED("friend_requested", "待加好友", AppColors.warning),
    READY_TO_CHAT("ready_to_chat", "可对话", AppColors.success),
    CHATTING("chatting", "沟通中", AppColors.success),
    CLOSED("closed", "已结束", AppColors.muted);

    companion object {
        fun of(key: String): SopStage? = entries.firstOrNull { it.key == key }
        fun label(key: String): String = of(key)?.label ?: key
        fun dot(key: String): Color = of(key)?.dot ?: AppColors.muted
    }
}
