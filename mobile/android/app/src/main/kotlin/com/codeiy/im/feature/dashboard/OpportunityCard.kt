package com.codeiy.im.feature.dashboard

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Group
import androidx.compose.material.icons.filled.Link
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.codeiy.im.model.Opportunity
import com.codeiy.im.model.Priority
import com.codeiy.im.ui.relativeTime
import com.codeiy.im.ui.theme.AppBadge
import com.codeiy.im.ui.theme.AppCard
import com.codeiy.im.ui.theme.AppColors
import com.codeiy.im.ui.theme.ConfidenceBadge
import com.codeiy.im.ui.theme.SopStage
import com.codeiy.im.ui.theme.TrustLevel
import com.codeiy.im.ui.theme.color
import kotlinx.serialization.json.JsonPrimitive

/** 商机卡片：信息层级对齐 Web opportunity-card.tsx 与 iOS OpportunityCardView。 */
@Composable
fun OpportunityCard(opportunity: Opportunity, onClick: () -> Unit, modifier: Modifier = Modifier) {
    val trust = TrustLevel.from(opportunity.trustScore)
    val linkStatus = (opportunity.linkVerification["status"] as? JsonPrimitive)?.content
    val showLinkRisk = opportunity.rawMessageLinks.isNotEmpty() &&
        (linkStatus == "unverified" || linkStatus == "verifying")

    AppCard(modifier = modifier.fillMaxWidth().clickable(onClick = onClick)) {
        Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
            Row(horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                Box(
                    Modifier.size(40.dp).background(opportunity.platform.color().copy(alpha = 0.15f), CircleShape),
                    contentAlignment = Alignment.Center,
                ) {
                    Text(
                        opportunity.contactName.take(1),
                        style = MaterialTheme.typography.titleMedium,
                        color = opportunity.platform.color(),
                    )
                }
                Column(Modifier.weight(1f), verticalArrangement = Arrangement.spacedBy(4.dp)) {
                    Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                        Text(
                            opportunity.contactName,
                            style = MaterialTheme.typography.titleMedium,
                            maxLines = 1,
                            overflow = TextOverflow.Ellipsis,
                            modifier = Modifier.weight(1f, fill = false),
                        )
                        if (opportunity.attentionRequired) {
                            Icon(Icons.Filled.Warning, contentDescription = "重大商机", tint = AppColors.destructive, modifier = Modifier.size(16.dp))
                        }
                    }
                    Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                        AppBadge(opportunity.platform.label, opportunity.platform.color())
                        if (opportunity.priority == Priority.HIGH || opportunity.priority == Priority.URGENT) {
                            AppBadge(opportunity.priority.label, opportunity.priority.color())
                        }
                        AppBadge(trust.label, trust.color())
                    }
                }
                ConfidenceBadge(opportunity.confidenceScore)
            }

            Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                if (opportunity.sourceType == "group") {
                    Icon(Icons.Filled.Group, contentDescription = null, tint = AppColors.muted, modifier = Modifier.size(14.dp))
                    Text(opportunity.groupName ?: "群消息", style = MaterialTheme.typography.labelSmall, color = AppColors.muted, maxLines = 1, overflow = TextOverflow.Ellipsis)
                } else {
                    Icon(Icons.Filled.Person, contentDescription = null, tint = AppColors.muted, modifier = Modifier.size(14.dp))
                    Text("私聊", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                }
                Text("·", color = AppColors.muted)
                Text(relativeTime(opportunity.createdAt), style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                Spacer(Modifier.weight(1f))
                AppBadge(opportunity.status.label, opportunity.status.color())
            }

            Text(
                opportunity.summary,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant,
                maxLines = 3,
                overflow = TextOverflow.Ellipsis,
            )

            if (showLinkRisk) {
                Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(4.dp)) {
                    Icon(Icons.Filled.Link, contentDescription = null, tint = AppColors.warning, modifier = Modifier.size(14.dp))
                    Text("含未核验链接，请先完成安全分析", style = MaterialTheme.typography.labelSmall, color = AppColors.warning)
                }
            }

            Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                Box(Modifier.size(6.dp).background(SopStage.dot(opportunity.sopStage), CircleShape))
                Text(SopStage.label(opportunity.sopStage), style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                opportunity.matchedKeywords.take(2).forEach { keyword ->
                    Text(keyword, style = MaterialTheme.typography.labelSmall, modifier = Modifier.background(MaterialTheme.colorScheme.surfaceVariant, CircleShape).padding(horizontal = 5.dp, vertical = 1.dp))
                }
                if (opportunity.matchedKeywords.size > 2) {
                    Text("+${opportunity.matchedKeywords.size - 2}", style = MaterialTheme.typography.labelSmall, color = AppColors.muted)
                }
            }
        }
    }
}
