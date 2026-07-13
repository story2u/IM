import SwiftUI

@main
struct OpportunityRadarApp: App {
    @State private var session = SessionStore()

    var body: some Scene {
        WindowGroup {
            appRoot
                .environment(session)
        }
    }

    @ViewBuilder
    private var appRoot: some View {
#if DEBUG
        if let screen = ProcessInfo.processInfo.demoScreen {
            DemoMobileRootView(screen: screen)
        } else {
            RootView()
        }
#else
        RootView()
#endif
    }
}

#if DEBUG
private extension ProcessInfo {
    var demoScreen: String? {
        guard let index = arguments.firstIndex(of: "-demo-screen"), arguments.indices.contains(index + 1) else {
            return nil
        }
        return arguments[index + 1]
    }
}
#endif

struct RootView: View {
    @Environment(SessionStore.self) private var session

    var body: some View {
        switch session.state {
        case .restoring:
            ProgressView("正在恢复会话…")
                .task { await session.restore() }
        case .restoreFailed(let message):
            ContentUnavailableView {
                Label("会话恢复失败", systemImage: "wifi.exclamationmark")
            } description: {
                Text(message)
            } actions: {
                Button("重试") { Task { await session.restore() } }
                Button("退出登录", role: .destructive) { session.logout() }
            }
        case .loggedOut:
            LoginView()
        case .active:
            InboxView()
        }
    }
}
