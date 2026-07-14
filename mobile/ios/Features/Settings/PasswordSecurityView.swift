import SwiftUI

struct PasswordSecurityView: View {
    @Environment(SessionStore.self) private var session
    @State private var currentPassword = ""
    @State private var newPassword = ""
    @State private var confirmPassword = ""
    @State private var isLoading = false
    @State private var showReset = false
    @State private var errorMessage: String?

    var body: some View {
        Form {
            Section {
                LabeledContent("登录邮箱", value: session.currentUser?.email ?? "—")
            }

            if session.currentUser?.hasPassword == true {
                Section("修改密码") {
                    SecureField("当前密码", text: $currentPassword)
                        .textContentType(.password)
                    SecureField("新密码", text: $newPassword)
                        .textContentType(.newPassword)
                    SecureField("确认新密码", text: $confirmPassword)
                        .textContentType(.newPassword)
                    Button {
                        Task { await changePassword() }
                    } label: {
                        if isLoading { ProgressView() } else { Text("修改密码") }
                    }
                    .disabled(!canChange || isLoading)
                } footer: {
                    Text("新密码至少 10 个字符。修改后所有设备需要重新登录。")
                }
            } else {
                Section {
                    Text("当前账户通过第三方登录。设置密码前需要验证登录邮箱。")
                    Button("验证邮箱并设置密码") { showReset = true }
                }
            }
        }
        .navigationTitle("账户安全")
        .sheet(isPresented: $showReset) {
            PasswordResetView(api: session.api, initialEmail: session.currentUser?.email ?? "")
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

    private var canChange: Bool {
        !currentPassword.isEmpty && newPassword.count >= 10 && newPassword == confirmPassword
    }

    @MainActor
    private func changePassword() async {
        guard canChange else {
            errorMessage = "请确认两次输入的新密码一致"
            return
        }
        isLoading = true
        defer { isLoading = false }
        do {
            try await session.changePassword(
                currentPassword: currentPassword,
                newPassword: newPassword
            )
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}
