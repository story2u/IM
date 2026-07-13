import SwiftUI

/// 商机识别规则：关键词 Chip + 添加/删除 + AI 语义识别开关。保存失败回滚。
struct DetectionSettingsView: View {
    let model: SettingsViewModel
    @State private var keywords: [String]
    @State private var aiSemanticsEnabled: Bool
    @State private var newKeyword = ""
    @State private var isSaving = false
    @State private var errorMessage: String?

    init(model: SettingsViewModel, detection: DetectionSettings) {
        self.model = model
        _keywords = State(initialValue: detection.keywords)
        _aiSemanticsEnabled = State(initialValue: detection.aiSemanticsEnabled)
    }

    var body: some View {
        Form {
            Section {
                Toggle(String(localized: "detection.ai_semantics", defaultValue: "AI 语义识别"), isOn: $aiSemanticsEnabled)
            } footer: {
                Text(String(localized: "detection.ai_hint", defaultValue: "开启后除关键词外，AI 会理解语义识别潜在商机。"))
            }

            Section {
                if keywords.isEmpty {
                    Text(String(localized: "detection.no_keywords", defaultValue: "暂无关键词")).foregroundStyle(.secondary)
                } else {
                    ForEach(keywords, id: \.self) { keyword in
                        HStack {
                            Text(keyword)
                            Spacer()
                            Button {
                                keywords.removeAll { $0 == keyword }
                            } label: {
                                Image(systemName: "xmark.circle.fill").foregroundStyle(.secondary)
                            }
                            .accessibilityLabel(Text("删除 \(keyword)"))
                        }
                    }
                }
                HStack {
                    TextField(String(localized: "detection.add_placeholder", defaultValue: "添加关键词"), text: $newKeyword)
                        .autocorrectionDisabled()
                        .onSubmit(addKeyword)
                    Button(String(localized: "action.add", defaultValue: "添加"), action: addKeyword)
                        .disabled(newKeyword.trimmingCharacters(in: .whitespaces).isEmpty)
                }
            } header: {
                Text(String(localized: "detection.keywords", defaultValue: "关键词"))
            }

            if let errorMessage {
                Section { Label(errorMessage, systemImage: "exclamationmark.triangle").foregroundStyle(AppColors.destructive) }
            }
        }
        .navigationTitle(Text("settings.detection", bundle: .main))
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .confirmationAction) {
                Button(String(localized: "action.save", defaultValue: "保存"), action: save).disabled(isSaving)
            }
        }
    }

    private func addKeyword() {
        let trimmed = newKeyword.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty, !keywords.contains(trimmed) else { newKeyword = ""; return }
        keywords.append(trimmed)
        newKeyword = ""
    }

    private func save() {
        isSaving = true
        errorMessage = nil
        Task {
            do {
                try await model.saveDetection(keywords: keywords, aiSemanticsEnabled: aiSemanticsEnabled)
            } catch {
                errorMessage = error.localizedDescription
                // 回滚后同步本地编辑态，避免展示已被服务端拒绝的值。
                if let saved = model.bundle?.detection {
                    keywords = saved.keywords
                    aiSemanticsEnabled = saved.aiSemanticsEnabled
                }
            }
            isSaving = false
        }
    }
}
