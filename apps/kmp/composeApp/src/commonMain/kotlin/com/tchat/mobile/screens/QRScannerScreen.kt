package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

data class QRResult(
    val type: QRResultType,
    val data: Map<String, Any?>
)

enum class QRResultType {
    PAYMENT, CONTACT, MERCHANT, PRODUCT, URL
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun QRScannerScreen(
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var flashEnabled by remember { mutableStateOf(false) }
    var scanResult by remember { mutableStateOf<QRResult?>(null) }
    var showPayment by remember { mutableStateOf(false) }
    var paymentAmount by remember { mutableStateOf("") }

    // Mock QR scan results for demo (following web design)
    val mockScanResults = mapOf(
        "PromptPay" to QRResult(
            QRResultType.PAYMENT,
            mapOf(
                "method" to "PromptPay",
                "recipient" to "Somtam Vendor",
                "phone" to "+66 XX XXX XXXX",
                "merchantId" to "merchant_123",
                "amount" to null
            )
        ),
        "Product" to QRResult(
            QRResultType.PRODUCT,
            mapOf(
                "id" to "prod_123",
                "name" to "Pad Thai Goong",
                "price" to 45,
                "merchant" to "Bangkok Street Food"
            )
        ),
        "Contact" to QRResult(
            QRResultType.CONTACT,
            mapOf(
                "name" to "John Doe",
                "phone" to "+66 XX XXX XXXX"
            )
        )
    )

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("QR Scanner", fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(
                            Icons.Default.ArrowBack,
                            contentDescription = "Back",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                actions = {
                    IconButton(onClick = { flashEnabled = !flashEnabled }) {
                        Icon(
                            if (flashEnabled) Icons.Default.FlashOn else Icons.Default.FlashOff,
                            contentDescription = "Flash",
                            tint = if (flashEnabled) TchatColors.primary else TchatColors.onSurface
                        )
                    }
                    IconButton(onClick = { /* Gallery */ }) {
                        Icon(
                            Icons.Default.Image,
                            contentDescription = "Gallery",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface,
                    titleContentColor = TchatColors.onSurface
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = modifier
                .fillMaxSize()
                .padding(paddingValues)
                .background(TchatColors.background)
        ) {
            // Camera Preview Area (Mock)
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(400.dp)
                    .padding(TchatSpacing.md),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceDim)
            ) {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
                    ) {
                        // QR Scanner Frame
                        Box(
                            modifier = Modifier
                                .size(200.dp)
                                .border(
                                    width = 3.dp,
                                    color = TchatColors.primary,
                                    shape = RoundedCornerShape(16.dp)
                                ),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                Icons.Default.QrCode,
                                contentDescription = "QR Scanner",
                                modifier = Modifier.size(80.dp),
                                tint = TchatColors.primary.copy(alpha = 0.5f)
                            )
                        }

                        Text(
                            "Position QR code within the frame",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )

                        if (flashEnabled) {
                            Row(
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(
                                    Icons.Default.FlashOn,
                                    contentDescription = null,
                                    tint = TchatColors.primary,
                                    modifier = Modifier.size(16.dp)
                                )
                                Spacer(modifier = Modifier.width(4.dp))
                                Text(
                                    "Flash On",
                                    style = MaterialTheme.typography.bodySmall,
                                    color = TchatColors.primary
                                )
                            }
                        }
                    }
                }
            }

            // Quick Action Buttons
            Text(
                "Quick Actions",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Bold,
                color = TchatColors.onSurface,
                modifier = Modifier.padding(horizontal = TchatSpacing.md, vertical = TchatSpacing.sm)
            )

            LazyColumn(
                modifier = Modifier.weight(1f),
                contentPadding = PaddingValues(horizontal = TchatSpacing.md),
                verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                items(mockScanResults.entries.toList()) { (label, result) ->
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable {
                                scanResult = result
                                if (result.type == QRResultType.PAYMENT) {
                                    showPayment = true
                                }
                            },
                        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                    ) {
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(TchatSpacing.md),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Box(
                                modifier = Modifier
                                    .size(48.dp)
                                    .clip(CircleShape)
                                    .background(
                                        when (result.type) {
                                            QRResultType.PAYMENT -> TchatColors.primary
                                            QRResultType.PRODUCT -> TchatColors.success
                                            QRResultType.CONTACT -> TchatColors.warning
                                            QRResultType.MERCHANT -> TchatColors.primaryLight
                                            QRResultType.URL -> TchatColors.surfaceVariant
                                        }
                                    ),
                                contentAlignment = Alignment.Center
                            ) {
                                val icon = when (result.type) {
                                    QRResultType.PAYMENT -> Icons.Default.Payment
                                    QRResultType.PRODUCT -> Icons.Default.ShoppingCart
                                    QRResultType.CONTACT -> Icons.Default.Person
                                    QRResultType.MERCHANT -> Icons.Default.Store
                                    QRResultType.URL -> Icons.Default.Link
                                }
                                Icon(
                                    icon,
                                    contentDescription = null,
                                    tint = TchatColors.onPrimary,
                                    modifier = Modifier.size(24.dp)
                                )
                            }

                            Spacer(modifier = Modifier.width(TchatSpacing.md))

                            Column(modifier = Modifier.weight(1f)) {
                                Text(
                                    "Scan $label QR",
                                    style = MaterialTheme.typography.bodyLarge,
                                    fontWeight = FontWeight.Medium,
                                    color = TchatColors.onSurface
                                )
                                Text(
                                    when (result.type) {
                                        QRResultType.PAYMENT -> "PromptPay payment"
                                        QRResultType.PRODUCT -> "Product information"
                                        QRResultType.CONTACT -> "Contact details"
                                        QRResultType.MERCHANT -> "Merchant info"
                                        QRResultType.URL -> "Website link"
                                    },
                                    style = MaterialTheme.typography.bodyMedium,
                                    color = TchatColors.onSurfaceVariant
                                )
                            }

                            Icon(
                                Icons.Default.ChevronRight,
                                contentDescription = null,
                                tint = TchatColors.onSurfaceVariant
                            )
                        }
                    }
                }
            }
        }

        // Payment Dialog
        if (showPayment) {
            PaymentDialog(
                qrResult = scanResult,
                paymentAmount = paymentAmount,
                onAmountChange = { paymentAmount = it },
                onConfirm = {
                    showPayment = false
                    scanResult = null
                    paymentAmount = ""
                    // Handle payment
                },
                onDismiss = {
                    showPayment = false
                    scanResult = null
                    paymentAmount = ""
                }
            )
        }
    }
}

@Composable
private fun PaymentDialog(
    qrResult: QRResult?,
    paymentAmount: String,
    onAmountChange: (String) -> Unit,
    onConfirm: () -> Unit,
    onDismiss: () -> Unit
) {
    Dialog(onDismissRequest = onDismiss) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            shape = RoundedCornerShape(16.dp),
            colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
        ) {
            Column(
                modifier = Modifier.padding(TchatSpacing.lg),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
            ) {
                Icon(
                    Icons.Default.Payment,
                    contentDescription = null,
                    modifier = Modifier.size(48.dp),
                    tint = TchatColors.primary
                )

                Text(
                    "PromptPay Payment",
                    style = MaterialTheme.typography.headlineSmall,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface
                )

                qrResult?.data?.let { data ->
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        Text(
                            data["recipient"] as? String ?: "Unknown",
                            style = MaterialTheme.typography.titleMedium,
                            color = TchatColors.onSurface
                        )
                        Text(
                            data["phone"] as? String ?: "",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }

                TchatInput(
                    value = paymentAmount,
                    onValueChange = onAmountChange,
                    type = TchatInputType.Number,
                    placeholder = "Enter amount (฿)",
                    modifier = Modifier.fillMaxWidth()
                )

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    TchatButton(
                        onClick = onDismiss,
                        text = "Cancel",
                        variant = TchatButtonVariant.Secondary,
                        modifier = Modifier.weight(1f)
                    )
                    TchatButton(
                        onClick = onConfirm,
                        text = "Pay",
                        variant = TchatButtonVariant.Primary,
                        modifier = Modifier.weight(1f)
                    )
                }
            }
        }
    }
}