import SwiftUI

/// 语义徽标：胶囊底色 + 同色前景（对齐 Web badge 与 Android AppBadge）。
struct AppBadge: View {
    let text: String
    let color: Color

    var body: some View {
        Text(text)
            .font(.caption2.weight(.medium))
            .padding(.horizontal, 6)
            .padding(.vertical, 2)
            .background(color.opacity(0.15), in: Capsule())
            .foregroundStyle(color)
    }
}

/// 相关度紧凑分数展示（环形在手机上过小，用带边框的百分比胶囊，语义等价）。
struct ConfidenceBadge: View {
    let score: Double

    private var percent: Int { Int((score * 100).rounded()) }

    var body: some View {
        Text("\(percent)%")
            .font(.caption2.weight(.semibold).monospacedDigit())
            .padding(.horizontal, 6)
            .padding(.vertical, 2)
            .overlay(Capsule().stroke(AppColors.primary.opacity(0.4), lineWidth: 1))
            .foregroundStyle(AppColors.primary)
            .accessibilityLabel(Text("相关度 \(percent)%"))
    }
}

/// 统一卡片容器：圆角 + 背景 + 阴影（对齐 Web card 与 Android AppCard）。
struct AppCard<Content: View>: View {
    @ViewBuilder var content: Content

    var body: some View {
        content
            .padding(12)
            .background(Color(.secondarySystemGroupedBackground), in: RoundedRectangle(cornerRadius: 12))
    }
}
