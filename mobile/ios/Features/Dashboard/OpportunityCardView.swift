import SwiftUI

/// 商机卡片：信息层级对齐 Web opportunity-card.tsx。
struct OpportunityCardView: View {
    let opportunity: Opportunity

    private var trustLevel: TrustLevel { TrustLevel.from(score: opportunity.trustScore) }

    /// 含未核验/分析中链接时提示先做安全分析（对齐 Web 逻辑）。
    private var showsLinkRisk: Bool {
        guard !opportunity.rawMessageLinks.isEmpty else { return false }
        if case .string(let status)? = opportunity.linkVerification["status"] {
            return status == "unverified" || status == "verifying"
        }
        return false
    }

    var body: some View {
        AppCard {
            VStack(alignment: .leading, spacing: 8) {
                header
                sourceRow
                Text(opportunity.summary)
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
                    .lineLimit(3)
                if showsLinkRisk {
                    Label(
                        String(localized: "card.link_risk", defaultValue: "含未核验链接，请先完成安全分析"),
                        systemImage: "link.badge.plus"
                    )
                    .font(.caption)
                    .foregroundStyle(AppColors.warning)
                }
                footer
            }
        }
    }

    private var header: some View {
        HStack(alignment: .top, spacing: 10) {
            avatar
            VStack(alignment: .leading, spacing: 4) {
                HStack(spacing: 6) {
                    Text(opportunity.contactName)
                        .font(.headline)
                        .lineLimit(1)
                    if opportunity.attentionRequired {
                        Image(systemName: "exclamationmark.circle.fill")
                            .foregroundStyle(AppColors.destructive)
                            .accessibilityLabel(Text("重大商机"))
                    }
                }
                HStack(spacing: 6) {
                    AppBadge(text: opportunity.platform.label, color: opportunity.platform.semanticColor)
                    if opportunity.priority == .high || opportunity.priority == .urgent {
                        AppBadge(text: opportunity.priority.label, color: opportunity.priority.semanticColor)
                    }
                    AppBadge(text: trustLevel.label, color: trustLevel.semanticColor)
                }
            }
            Spacer(minLength: 4)
            ConfidenceBadge(score: opportunity.confidenceScore)
        }
    }

    private var avatar: some View {
        ZStack {
            Circle().fill(opportunity.platform.semanticColor.opacity(0.15))
            Text(String(opportunity.contactName.prefix(1)))
                .font(.headline)
                .foregroundStyle(opportunity.platform.semanticColor)
        }
        .frame(width: 40, height: 40)
    }

    private var sourceRow: some View {
        HStack(spacing: 6) {
            if opportunity.sourceType == "group" {
                Image(systemName: "person.3.fill").font(.caption2).foregroundStyle(.secondary)
                Text(opportunity.groupName ?? String(localized: "source.group", defaultValue: "群消息"))
                    .font(.caption)
                    .foregroundStyle(.secondary)
                    .lineLimit(1)
            } else {
                Image(systemName: "person.fill").font(.caption2).foregroundStyle(.secondary)
                Text(String(localized: "source.private_short", defaultValue: "私聊"))
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }
            Text("·").foregroundStyle(.secondary)
            Text(opportunity.createdAt, format: .relative(presentation: .named))
                .font(.caption)
                .foregroundStyle(.secondary)
            Spacer()
            AppBadge(text: opportunity.status.label, color: opportunity.status.semanticColor)
        }
    }

    private var footer: some View {
        HStack(spacing: 6) {
            HStack(spacing: 4) {
                Circle().fill(SopStage(rawValue: opportunity.sopStage)?.dotColor ?? AppColors.muted)
                    .frame(width: 6, height: 6)
                Text(SopStage.label(for: opportunity.sopStage))
                    .font(.caption2)
                    .foregroundStyle(.secondary)
            }
            if !opportunity.matchedKeywords.isEmpty {
                let shown = opportunity.matchedKeywords.prefix(2)
                ForEach(Array(shown), id: \.self) { keyword in
                    Text(keyword)
                        .font(.caption2)
                        .padding(.horizontal, 5)
                        .padding(.vertical, 1)
                        .background(Color(.tertiarySystemFill), in: Capsule())
                }
                if opportunity.matchedKeywords.count > shown.count {
                    Text("+\(opportunity.matchedKeywords.count - shown.count)")
                        .font(.caption2)
                        .foregroundStyle(.secondary)
                }
            }
            Spacer()
        }
    }
}
