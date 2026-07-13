import SwiftUI

/// SOP 流程阶段：键与顺序对齐后端 opportunity.sop_stage 与 Web sop.ts。
enum SopStage: String, CaseIterable, Identifiable {
    case detected
    case analyzing
    case verified
    case contactExtracted = "contact_extracted"
    case friendRequested = "friend_requested"
    case readyToChat = "ready_to_chat"
    case chatting
    case closed

    var id: String { rawValue }

    var label: String {
        switch self {
        case .detected: String(localized: "sop.detected", defaultValue: "已发现")
        case .analyzing: String(localized: "sop.analyzing", defaultValue: "AI 分析中")
        case .verified: String(localized: "sop.verified", defaultValue: "已核验")
        case .contactExtracted: String(localized: "sop.contact_extracted", defaultValue: "已提取联系方式")
        case .friendRequested: String(localized: "sop.friend_requested", defaultValue: "待加好友")
        case .readyToChat: String(localized: "sop.ready_to_chat", defaultValue: "可对话")
        case .chatting: String(localized: "sop.chatting", defaultValue: "沟通中")
        case .closed: String(localized: "sop.closed", defaultValue: "已结束")
        }
    }

    var dotColor: Color {
        switch self {
        case .detected, .closed: AppColors.muted
        case .analyzing, .verified, .contactExtracted: AppColors.primary
        case .friendRequested: AppColors.warning
        case .readyToChat, .chatting: AppColors.success
        }
    }

    /// 后端可能返回未知阶段，展示时容错。
    static func label(for raw: String) -> String {
        SopStage(rawValue: raw)?.label ?? raw
    }
}
