import SwiftUI

/// 通知偏好：4 个开关。推送通道未开发时明确标注"将在推送服务启用后生效"。
struct NotificationSettingsView: View {
    let model: SettingsViewModel
    let pushAvailable: Bool
    @State private var prefs: NotificationSettings
    @State private var errorMessage: String?

    init(model: SettingsViewModel, notifications: NotificationSettings, pushAvailable: Bool) {
        self.model = model
        self.pushAvailable = pushAvailable
        _prefs = State(initialValue: notifications)
    }

    var body: some View {
        Form {
            if !pushAvailable {
                Section {
                    Label(
                        String(localized: "notifications.push_unavailable", defaultValue: "推送服务尚未开放，偏好会保存，将在启用后生效。"),
                        systemImage: "info.circle"
                    )
                    .font(.footnote)
                    .foregroundStyle(.secondary)
                }
            }
            Section {
                Toggle(String(localized: "notifications.new_opportunity", defaultValue: "新商机提醒"), isOn: $prefs.newOpportunityEnabled)
                Toggle(String(localized: "notifications.ai_replied", defaultValue: "AI 已回复通知"), isOn: $prefs.aiRepliedEnabled)
                Toggle(String(localized: "notifications.daily_digest", defaultValue: "每日商机摘要"), isOn: $prefs.dailyDigestEnabled)
                Toggle(String(localized: "notifications.urgent_only", defaultValue: "仅紧急商机"), isOn: $prefs.urgentOnly)
            }
            if let errorMessage {
                Section { Label(errorMessage, systemImage: "exclamationmark.triangle").foregroundStyle(AppColors.destructive) }
            }
        }
        .navigationTitle(Text("settings.notifications", bundle: .main))
        .navigationBarTitleDisplayMode(.inline)
        // 开关即时保存；失败回滚到服务端值。
        .onChange(of: prefs) { _, newValue in save(newValue) }
    }

    private func save(_ newValue: NotificationSettings) {
        errorMessage = nil
        Task {
            do {
                try await model.saveNotifications(newValue)
            } catch {
                errorMessage = error.localizedDescription
                if let saved = model.bundle?.notifications { prefs = saved }
            }
        }
    }
}
