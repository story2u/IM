#if DEBUG
import SwiftUI

struct DemoMobileRootView: View {
    let screen: String

    var body: some View {
        NavigationStack {
            Group {
                switch screen {
                case "login": DemoLogin()
                case "opportunity-detail": DemoDetail()
                case "settings": DemoSettings()
                default: DemoDashboard()
                }
            }
            .preferredColorScheme(.light)
        }
    }
}

private struct DemoLogin: View {
    var body: some View { VStack(spacing: 24) { Spacer(); Image(systemName: "dot.radiowaves.left.and.right").font(.system(size: 64)).foregroundStyle(.indigo); Text("商机雷达").font(.largeTitle.bold()); Text("登录后发现聊天中的潜在商机").foregroundStyle(.secondary); Button("使用 Google 账号继续") {}.buttonStyle(.borderedProminent); Button("使用 Apple 账号继续") {}.buttonStyle(.bordered); Spacer() }.padding(28).navigationTitle("登录或注册") }
}

private struct DemoDashboard: View {
    var body: some View { List { Section { Label("Pi Agent 发现 2 条重大商机", systemImage: "exclamationmark.circle.fill").foregroundStyle(.orange) }; Section("待处理") { DemoRow(name: "林远（演示）", text: "采购 50 套设备，下周安排演示", score: "96%", channel: "Telegram"); DemoRow(name: "周屿（演示）", text: "寻找 API 服务商，需要 CRM 与 SLA", score: "93%", channel: "Telegram"); DemoRow(name: "顾言（演示）", text: "年度续约并增购 30 个席位", score: "91%", channel: "企业微信") } }.navigationTitle("商机收件箱") }
}

private struct DemoRow: View { let name: String; let text: String; let score: String; let channel: String; var body: some View { VStack(alignment: .leading, spacing: 7) { HStack { Text(name).font(.headline); Spacer(); Text(score).foregroundStyle(.indigo).bold() }; Text(text).font(.subheadline).foregroundStyle(.secondary); Text(channel).font(.caption).foregroundStyle(.blue) }.padding(.vertical, 4) } }

private struct DemoDetail: View {
    var body: some View { List { Section("原始需求") { Text("我们团队想采购 50 套设备，下周能安排演示吗？预算已确认。"); LabeledContent("相关度", value: "96%"); LabeledContent("可信度", value: "91") }; Section("Pi Agent") { Label("未包含外部链接", systemImage: "checkmark.shield"); Label("procurement@example.com", systemImage: "envelope"); Text("建议确认演示时间和设备规格，执行前需要人工批准。") }; Section("回复草稿") { Text("可以安排演示，请问下周二或周三哪天方便？"); Button("发送（演示禁用）") {}.disabled(true) } }.navigationTitle("商机详情") }
}

private struct DemoSettings: View {
    var body: some View { List { Section("订阅") { Label("Pro · 本月 18 / 500 次分析", systemImage: "creditcard") }; Section("连接") { Label("Telegram 普通账号 · 已连接", systemImage: "paperplane"); Label("企业微信 · 已连接", systemImage: "message") }; Section("工作时间") { LabeledContent("周一至周五", value: "09:00–18:00") }; Section("安全") { Label("外部动作人工审批", systemImage: "person.badge.shield.checkmark") } }.navigationTitle("设置中心") }
}
#endif
