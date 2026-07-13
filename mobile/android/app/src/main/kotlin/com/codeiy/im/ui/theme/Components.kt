package com.codeiy.im.ui.theme

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp

/** 语义徽标：胶囊底色 + 同色前景（对齐 iOS AppBadge / Web badge）。 */
@Composable
fun AppBadge(text: String, color: Color) {
    Surface(color = color.copy(alpha = 0.15f), shape = CircleShape) {
        Text(
            text,
            style = MaterialTheme.typography.labelSmall,
            color = color,
            modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
        )
    }
}

/** 统一卡片容器（对齐 iOS AppCard / Web card）。 */
@Composable
fun AppCard(modifier: Modifier = Modifier, content: @Composable () -> Unit) {
    Surface(
        modifier = modifier,
        shape = RoundedCornerShape(12.dp),
        color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.4f),
        tonalElevation = 1.dp,
    ) {
        Box(Modifier.padding(12.dp)) { content() }
    }
}

/** 相关度紧凑分数（环形在手机过小，用带边框百分比，语义等价）。 */
@Composable
fun ConfidenceBadge(score: Double) {
    val percent = (score * 100).toInt()
    Text(
        "$percent%",
        style = MaterialTheme.typography.labelSmall,
        color = MaterialTheme.colorScheme.primary,
    )
}
