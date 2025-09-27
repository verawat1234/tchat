package com.tchat.mobile.components

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TchatTopBar(
    title: String,
    onWebClick: () -> Unit,
    onMoreClick: () -> Unit,
    showWebButton: Boolean = true,
    showMoreButton: Boolean = true,
    actions: @Composable RowScope.() -> Unit = {},
    modifier: Modifier = Modifier
) {
    TopAppBar(
        title = {
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = title,
                    style = MaterialTheme.typography.titleLarge,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onPrimary
                )

                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    // Quick Web Access Button
                    if (showWebButton) {
                        TchatWebButton(onClick = onWebClick)
                    }

                    // More Button (+ icon)
                    if (showMoreButton) {
                        TchatMoreButton(onClick = onMoreClick)
                    }

                    // Custom actions
                    actions()
                }
            }
        },
        colors = TopAppBarDefaults.topAppBarColors(
            containerColor = TchatColors.primary,
            titleContentColor = TchatColors.onPrimary
        ),
        modifier = modifier
    )
}

@Composable
private fun TchatWebButton(
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    OutlinedButton(
        onClick = onClick,
        colors = ButtonDefaults.outlinedButtonColors(
            contentColor = TchatColors.onPrimary
        ),
        border = ButtonDefaults.outlinedButtonBorder.copy(
            brush = androidx.compose.ui.graphics.SolidColor(TchatColors.onPrimary)
        ),
        modifier = modifier.height(32.dp)
    ) {
        Icon(
            Icons.Default.Language,
            contentDescription = "Web Browser",
            modifier = Modifier.size(16.dp),
            tint = TchatColors.onPrimary
        )
        Spacer(modifier = Modifier.width(4.dp))
        Text(
            text = "Web",
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.onPrimary
        )
    }
}

@Composable
private fun TchatMoreButton(
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    OutlinedButton(
        onClick = onClick,
        colors = ButtonDefaults.outlinedButtonColors(
            contentColor = TchatColors.onPrimary
        ),
        border = ButtonDefaults.outlinedButtonBorder.copy(
            brush = androidx.compose.ui.graphics.SolidColor(TchatColors.onPrimary)
        ),
        modifier = modifier.height(32.dp)
    ) {
        Icon(
            Icons.Default.Add,
            contentDescription = "More Options",
            modifier = Modifier.size(16.dp),
            tint = TchatColors.onPrimary
        )
        Spacer(modifier = Modifier.width(4.dp))
        Text(
            text = "More",
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.onPrimary
        )
    }
}