import SwiftUI

struct PasswordResetView: View {
    let api: APIClient
    let initialEmail: String

    @Environment(\.dismiss) private var dismiss
    @State private var email: String
    @State private var code = ""
    @State private var newPassword = ""
    @State private var confirmPassword = ""
    @State private var requested = false
    @State private var isLoading = false
    @State private var message: String?
    @State private var errorMessage: String?

    init(api: APIClient, initialEmail: String = "") {
        self.api = api
        self.initialEmail = initialEmail
        _email = State(initialValue: initialEmail)
    }

    var body: some View {
        NavigationStack {
            Form {
                Section {
                    TextField("邮箱", text: $email)
                        .textContentType(.emailAddress)
                        .keyboardType(.emailAddress)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                        .disabled(requested)
                } footer: {
                    Text("无论账户是否存在，系统都会返回相同提示。")
                }

                if requested {
                    Section("邮件验证") {
                        TextField("10 位验证码", text: $code)
                            .textContentType(.oneTimeCode)
                            .textInputAutocapitalization(.characters)
                            .autocorrectionDisabled()
                            .onChange(of: code) { _, value in
                                code = String(value.uppercased().filter { !$0.isWhitespace }.prefix(10))
                            }
                        SecureField("新密码", text: $newPassword)
                            .textContentType(.newPassword)
                        SecureField("确认新密码", text: $confirmPassword)
                            .textContentType(.newPassword)
                    } footer: {
                        Text("新密码至少 10 个字符。")
                    }
                }

                if let message {
                    Section { Text(message).foregroundStyle(.secondary) }
                }

                Section {
                    Button {
                        Task { await submit() }
                    } label: {
                        HStack {
                            Spacer()
                            if isLoading { ProgressView() }
                            Text(requested ? "确认重置" : "发送重置邮件")
                            Spacer()
                        }
                    }
                    .disabled(!canSubmit || isLoading)
                }
            }
            .navigationTitle("重置密码")
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("关闭") { dismiss() }
                }
            }
            .alert("操作失败", isPresented: .init(
                get: { errorMessage != nil },
                set: { if !$0 { errorMessage = nil } }
            )) {
                Button("好", role: .cancel) {}
            } message: {
                Text(errorMessage ?? "")
            }
        }
    }

    private var canSubmit: Bool {
        let normalizedEmail = email.trimmingCharacters(in: .whitespacesAndNewlines)
        guard normalizedEmail.contains("@"), normalizedEmail.count <= 320 else { return false }
        if !requested { return true }
        return code.count == 10 && newPassword.count >= 10 && newPassword == confirmPassword
    }

    @MainActor
    private func submit() async {
        guard canSubmit else {
            if requested && newPassword != confirmPassword {
                errorMessage = "两次输入的新密码不一致"
            }
            return
        }
        isLoading = true
        defer { isLoading = false }
        do {
            if requested {
                _ = try await api.confirmPasswordReset(
                    email: email.trimmingCharacters(in: .whitespacesAndNewlines),
                    code: code,
                    newPassword: newPassword
                )
                dismiss()
            } else {
                let result = try await api.requestPasswordReset(
                    email: email.trimmingCharacters(in: .whitespacesAndNewlines)
                )
                message = result.message
                requested = true
            }
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}
