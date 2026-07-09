from typing import TypedDict

from langgraph.graph import END, StateGraph


class ReplyState(TypedDict):
    opportunity_id: str
    draft: str
    need_human: bool


RISKY_WORDS = {"保证", "最低价", "合同已确认", "一定交付", "100%"}


async def policy_check_node(state: ReplyState) -> ReplyState:
    state["need_human"] = any(word in state["draft"] for word in RISKY_WORDS)
    return state


def build_policy_graph():
    graph = StateGraph(ReplyState)
    graph.add_node("policy_check", policy_check_node)
    graph.set_entry_point("policy_check")
    graph.add_conditional_edges(
        "policy_check",
        lambda state: "human" if state["need_human"] else "send",
        {"human": END, "send": END},
    )
    return graph.compile()
