import SwiftUI

/// 跨端语义色（对齐 Web theme 与 Android ui/theme/Color.kt）。
/// 不追求像素一致，只保证语义一致：同一状态在三端用同一族颜色。
enum AppColors {
    static let primary = Color.accentColor
    static let success = Color.green
    static let warning = Color.orange
    static let destructive = Color.red
    static let muted = Color.secondary
    static let telegram = Color.blue
    static let wecom = Color.green
    static let ai = Color.purple
}

// MARK: - 语义映射（优先级/状态/可信度）

extension Priority {
    /// 紧急=destructive，高=warning，中=primary，低=muted。
    var semanticColor: Color {
        switch self {
        case .urgent: AppColors.destructive
        case .high: AppColors.warning
        case .normal: AppColors.primary
        case .low, .unknown: AppColors.muted
        }
    }
}

extension FrontendOpportunityStatus {
    var semanticColor: Color {
        switch self {
        case .pending: AppColors.warning
        case .replied: AppColors.success
        case .ignored, .unknown: AppColors.muted
        }
    }
}

extension IMChannel {
    var semanticColor: Color {
        switch self {
        case .telegram: AppColors.telegram
        case .wecom: AppColors.wecom
        case .unknown: AppColors.muted
        }
    }
}

/// 可信度等级：与后端 dashboard trust_levels 及 Web sop.ts 边界一致。
enum TrustLevel: String, CaseIterable, Identifiable {
    case trusted, unverified, suspicious, risky

    var id: String { rawValue }

    /// 后端存的是 0-100 分数；这里用同一边界反推等级用于展示。
    static func from(score: Int) -> TrustLevel {
        switch score {
        case 80...: .trusted
        case 60..<80: .unverified
        case 40..<60: .suspicious
        default: .risky
        }
    }

    var label: String {
        switch self {
        case .trusted: String(localized: "trust.trusted", defaultValue: "安全可信")
        case .unverified: String(localized: "trust.unverified", defaultValue: "待核验")
        case .suspicious: String(localized: "trust.suspicious", defaultValue: "可疑")
        case .risky: String(localized: "trust.risky", defaultValue: "高风险")
        }
    }

    var semanticColor: Color {
        switch self {
        case .trusted: AppColors.success
        case .unverified: AppColors.muted
        case .suspicious: AppColors.warning
        case .risky: AppColors.destructive
        }
    }
}
