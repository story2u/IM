import Observation
import SwiftUI

@MainActor
@Observable
final class TelegramSettingsModel {
    private let api: APIClient
    var health: TelegramConnectionHealth?
    var connections: [TelegramConnectionDTO] = []
    var isLoading = false
    var errorMessage: String?

    init(api: APIClient) { self.api = api }

    func load() async {
        isLoading = true
        defer { isLoading = false }
        do {
            async let health = api.telegramHealth()
            async let connections = api.telegramConnections()
            self.health = try await health
            self.connections = try await connections
            errorMessage = nil
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func toggle(_ connection: TelegramConnectionDTO) async {
        do {
            let updated = try await api.setTelegramConnectionEnabled(id: connection.id, enabled: !connection.enabled)
            if let index = connections.firstIndex(where: { $0.id == updated.id }) {
                connections[index] = updated
            }
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

/// Telegram 连接：真实读取 health + connections + sources，可停用/启用。
/// 连接建立（Bot/Business/QR 三卡握手 + 深链 + 轮询）为后续迭代，未配置能力如实标注。
struct TelegramSettingsView: View {
    @Environment(SessionStore.self) private var session
    @State private var model: TelegramSettingsModel?

    var body: some View {
        Group {
            if let model {
                content(model)
            } else {
                ProgressView()
            }
        }
        .navigationTitle(Text("settings.telegram", bundle: .main))
        .navigationBarTitleDisplayMode(.inline)
        .task {
            if model == nil {
                let vm = TelegramSettingsModel(api: session.api)
                model = vm
                await vm.load()
            }
        }
    }

    private func content(_ model: TelegramSettingsModel) -> some View {
        List {
            if let health = model.health {
                Section(String(localized: "telegram.status", defaultValue: "服务状态")) {
                    LabeledContent(String(localized: "telegram.bot", defaultValue: "Bot"), value: health.botConfigured
                        ? (health.botUsername.map { "@\($0)" } ?? String(localized: "telegram.configured", defaultValue: "已配置"))
                        : String(localized: "telegram.admin_unconfigured", defaultValue: "管理员尚未配置"))
                    capabilityRow(String(localized: "telegram.business", defaultValue: "Business 私聊"), available: health.businessAvailable)
                    capabilityRow(String(localized: "telegram.qr", defaultValue: "普通账号 QR"), available: health.mtprotoQrAvailable)
                }
            }

            if model.connections.isEmpty && !model.isLoading {
                Section {
                    ContentUnavailableView(
                        String(localized: "telegram.no_connections", defaultValue: "尚无连接"),
                        systemImage: "paperplane",
                        description: Text(String(localized: "telegram.connect_hint", defaultValue: "连接建立向导即将上线；当前可在 Web 端完成绑定后在此管理。"))
                    )
                }
            }

            ForEach(model.connections) { connection in
                Section(connection.label) {
                    Toggle(String(localized: "telegram.enabled", defaultValue: "启用"), isOn: Binding(
                        get: { connection.enabled },
                        set: { _ in Task { await model.toggle(connection) } }
                    ))
                    LabeledContent(String(localized: "telegram.conn_status", defaultValue: "状态"), value: connection.status)
                    if let error = connection.lastError {
                        Label(error, systemImage: "exclamationmark.triangle").font(.caption).foregroundStyle(AppColors.warning)
                    }
                    ForEach(connection.sources) { source in
                        HStack {
                            Image(systemName: source.sourceType == "private" ? "person" : "person.3")
                                .foregroundStyle(.secondary)
                            VStack(alignment: .leading) {
                                Text(source.displayName)
                                if source.quotaPaused, let reason = source.quotaReason {
                                    Text(reason).font(.caption2).foregroundStyle(AppColors.warning)
                                }
                            }
                            Spacer()
                            if !source.enabled { Text(String(localized: "telegram.disabled", defaultValue: "已停用")).font(.caption).foregroundStyle(.secondary) }
                        }
                    }
                }
            }

            if let error = model.errorMessage {
                Section { Label(error, systemImage: "exclamationmark.triangle").foregroundStyle(AppColors.destructive) }
            }
        }
        .refreshable { await model.load() }
    }

    private func capabilityRow(_ title: String, available: Bool) -> some View {
        LabeledContent(title, value: available
            ? String(localized: "telegram.available", defaultValue: "可用")
            : String(localized: "telegram.admin_unconfigured", defaultValue: "管理员尚未配置"))
    }
}
