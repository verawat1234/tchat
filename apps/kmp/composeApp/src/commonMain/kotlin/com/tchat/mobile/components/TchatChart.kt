package com.tchat.mobile.components

import androidx.compose.animation.core.*
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.*
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.designsystem.TchatColors
import kotlin.math.*

/**
 * TchatChart - Basic chart components
 *
 * Features:
 * - Line, Bar, Pie chart types with smooth animations
 * - Interactive data points with hover/click events
 * - Legend and axis labels with customization
 * - Responsive sizing and touch interaction
 */

data class ChartData(
    val value: Float,
    val label: String,
    val color: Color = TchatColors.primary
)

data class LineChartPoint(
    val x: Float,
    val y: Float,
    val label: String = ""
)

enum class ChartType {
    LINE,
    BAR,
    PIE,
    DOUGHNUT
}

@Composable
fun TchatChart(
    data: List<ChartData>,
    type: ChartType = ChartType.BAR,
    modifier: Modifier = Modifier,
    title: String? = null,
    showLegend: Boolean = true,
    showLabels: Boolean = true,
    showValues: Boolean = false,
    animated: Boolean = true,
    interactive: Boolean = true,
    onDataPointClick: ((ChartData, Int) -> Unit)? = null
) {
    Column(
        modifier = modifier.fillMaxWidth(),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        // Title
        title?.let {
            Text(
                text = it,
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Medium,
                color = TchatColors.onSurface
            )
        }

        // Chart content
        when (type) {
            ChartType.LINE -> {
                // Convert ChartData to LineChartPoint for line chart
                val linePoints = data.mapIndexed { index, chartData ->
                    LineChartPoint(index.toFloat(), chartData.value, chartData.label)
                }
                LineChart(
                    points = linePoints,
                    color = data.firstOrNull()?.color ?: TchatColors.primary,
                    showLabels = showLabels,
                    showValues = showValues,
                    animated = animated,
                    interactive = interactive,
                    onPointClick = if (onDataPointClick != null) { point, index ->
                        data.getOrNull(index)?.let { onDataPointClick(it, index) }
                    } else null
                )
            }
            ChartType.BAR -> {
                BarChart(
                    data = data,
                    showLabels = showLabels,
                    showValues = showValues,
                    animated = animated,
                    interactive = interactive,
                    onBarClick = onDataPointClick
                )
            }
            ChartType.PIE -> {
                PieChart(
                    data = data,
                    showLabels = showLabels,
                    showValues = showValues,
                    animated = animated,
                    interactive = interactive,
                    onSliceClick = onDataPointClick
                )
            }
            ChartType.DOUGHNUT -> {
                DoughnutChart(
                    data = data,
                    showLabels = showLabels,
                    showValues = showValues,
                    animated = animated,
                    interactive = interactive,
                    onSliceClick = onDataPointClick
                )
            }
        }

        // Legend
        if (showLegend) {
            ChartLegend(data = data)
        }
    }
}

@Composable
private fun LineChart(
    points: List<LineChartPoint>,
    color: Color,
    showLabels: Boolean,
    showValues: Boolean,
    animated: Boolean,
    interactive: Boolean,
    onPointClick: ((LineChartPoint, Int) -> Unit)?
) {
    if (points.isEmpty()) return

    val animationProgress by animateFloatAsState(
        targetValue = if (animated) 1f else 1f,
        animationSpec = if (animated) {
            tween(durationMillis = 1000, easing = FastOutSlowInEasing)
        } else {
            snap()
        },
        label = "line_chart_animation"
    )

    Column {
        // Chart area
        Canvas(
            modifier = Modifier
                .fillMaxWidth()
                .height(200.dp)
                .semantics {
                    contentDescription = "Line chart with ${points.size} data points"
                }
        ) {
            drawLineChart(
                points = points,
                color = color,
                animationProgress = animationProgress,
                size = size
            )
        }

        // Labels
        if (showLabels) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                points.forEach { point ->
                    Text(
                        text = point.label,
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant,
                        textAlign = TextAlign.Center,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis,
                        modifier = Modifier.weight(1f)
                    )
                }
            }
        }
    }
}

@Composable
private fun BarChart(
    data: List<ChartData>,
    showLabels: Boolean,
    showValues: Boolean,
    animated: Boolean,
    interactive: Boolean,
    onBarClick: ((ChartData, Int) -> Unit)?
) {
    if (data.isEmpty()) return

    val maxValue = data.maxOfOrNull { it.value } ?: 1f

    Column {
        // Chart bars
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .height(200.dp),
            horizontalArrangement = Arrangement.spacedBy(8.dp),
            verticalAlignment = Alignment.Bottom
        ) {
            data.forEachIndexed { index, chartData ->
                val animatedHeight by animateFloatAsState(
                    targetValue = if (animated) chartData.value / maxValue else chartData.value / maxValue,
                    animationSpec = if (animated) {
                        tween(durationMillis = 800, delayMillis = index * 100, easing = FastOutSlowInEasing)
                    } else {
                        snap()
                    },
                    label = "bar_height_$index"
                )

                Column(
                    modifier = Modifier.weight(1f),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Bottom
                ) {
                    // Value label
                    if (showValues) {
                        Text(
                            text = chartData.value.toInt().toString(),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant,
                            fontSize = 10.sp
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                    }

                    // Bar
                    Box(
                        modifier = Modifier
                            .width(32.dp)
                            .fillMaxHeight(animatedHeight)
                            .clip(RoundedCornerShape(topStart = 4.dp, topEnd = 4.dp))
                            .background(chartData.color)
                            .run {
                                if (interactive && onBarClick != null) {
                                    clickable { onBarClick(chartData, index) }
                                } else this
                            }
                            .semantics {
                                contentDescription = "${chartData.label}: ${chartData.value}"
                                role = Role.Button
                            }
                    )
                }
            }
        }

        // Labels
        if (showLabels) {
            Spacer(modifier = Modifier.height(8.dp))
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                data.forEach { chartData ->
                    Text(
                        text = chartData.label,
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant,
                        textAlign = TextAlign.Center,
                        maxLines = 2,
                        overflow = TextOverflow.Ellipsis,
                        modifier = Modifier.weight(1f)
                    )
                }
            }
        }
    }
}

@Composable
private fun PieChart(
    data: List<ChartData>,
    showLabels: Boolean,
    showValues: Boolean,
    animated: Boolean,
    interactive: Boolean,
    onSliceClick: ((ChartData, Int) -> Unit)?
) {
    if (data.isEmpty()) return

    val total = data.sumOf { it.value.toDouble() }.toFloat()
    var startAngle = -90f

    val animationProgress by animateFloatAsState(
        targetValue = if (animated) 1f else 1f,
        animationSpec = if (animated) {
            tween(durationMillis = 1000, easing = FastOutSlowInEasing)
        } else {
            snap()
        },
        label = "pie_chart_animation"
    )

    Box(
        modifier = Modifier
            .size(200.dp)
            .semantics {
                contentDescription = "Pie chart with ${data.size} segments"
            },
        contentAlignment = Alignment.Center
    ) {
        Canvas(
            modifier = Modifier.fillMaxSize()
        ) {
            data.forEachIndexed { index, chartData ->
                val sweepAngle = (360f * chartData.value / total) * animationProgress

                drawArc(
                    color = chartData.color,
                    startAngle = startAngle,
                    sweepAngle = sweepAngle,
                    useCenter = true,
                    size = Size(size.width * 0.8f, size.height * 0.8f),
                    topLeft = Offset(size.width * 0.1f, size.height * 0.1f)
                )

                startAngle += sweepAngle
            }
        }

        // Center text (total value)
        if (showValues) {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = total.toInt().toString(),
                    style = MaterialTheme.typography.headlineMedium,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface
                )
                Text(
                    text = "Total",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
private fun DoughnutChart(
    data: List<ChartData>,
    showLabels: Boolean,
    showValues: Boolean,
    animated: Boolean,
    interactive: Boolean,
    onSliceClick: ((ChartData, Int) -> Unit)?
) {
    if (data.isEmpty()) return

    val total = data.sumOf { it.value.toDouble() }.toFloat()
    var startAngle = -90f
    val strokeWidth = with(LocalDensity.current) { 40.dp.toPx() }

    val animationProgress by animateFloatAsState(
        targetValue = if (animated) 1f else 1f,
        animationSpec = if (animated) {
            tween(durationMillis = 1000, easing = FastOutSlowInEasing)
        } else {
            snap()
        },
        label = "doughnut_chart_animation"
    )

    Box(
        modifier = Modifier
            .size(200.dp)
            .semantics {
                contentDescription = "Doughnut chart with ${data.size} segments"
            },
        contentAlignment = Alignment.Center
    ) {
        Canvas(
            modifier = Modifier.fillMaxSize()
        ) {
            data.forEachIndexed { index, chartData ->
                val sweepAngle = (360f * chartData.value / total) * animationProgress

                drawArc(
                    color = chartData.color,
                    startAngle = startAngle,
                    sweepAngle = sweepAngle,
                    useCenter = false,
                    style = Stroke(width = strokeWidth, cap = StrokeCap.Round),
                    size = Size(size.width - strokeWidth, size.height - strokeWidth),
                    topLeft = Offset(strokeWidth / 2, strokeWidth / 2)
                )

                startAngle += sweepAngle
            }
        }

        // Center content
        if (showValues) {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = total.toInt().toString(),
                    style = MaterialTheme.typography.headlineMedium,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface
                )
                Text(
                    text = "Total",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
private fun ChartLegend(data: List<ChartData>) {
    LazyRow(
        horizontalArrangement = Arrangement.spacedBy(16.dp),
        contentPadding = PaddingValues(horizontal = 16.dp)
    ) {
        items(data) { chartData ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Box(
                    modifier = Modifier
                        .size(12.dp)
                        .clip(CircleShape)
                        .background(chartData.color)
                )

                Text(
                    text = chartData.label,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurface
                )

                Text(
                    text = "(${chartData.value.toInt()})",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

private fun DrawScope.drawLineChart(
    points: List<LineChartPoint>,
    color: Color,
    animationProgress: Float,
    size: Size
) {
    if (points.size < 2) return

    val maxX = points.maxOfOrNull { it.x } ?: 1f
    val minX = points.minOfOrNull { it.x } ?: 0f
    val maxY = points.maxOfOrNull { it.y } ?: 1f
    val minY = points.minOfOrNull { it.y } ?: 0f

    val padding = 32.dp.toPx()
    val chartWidth = size.width - 2 * padding
    val chartHeight = size.height - 2 * padding

    // Convert points to screen coordinates
    val screenPoints = points.map { point ->
        Offset(
            x = padding + (point.x - minX) / (maxX - minX) * chartWidth,
            y = size.height - padding - (point.y - minY) / (maxY - minY) * chartHeight
        )
    }

    // Draw animated line
    val animatedPointsCount = (screenPoints.size * animationProgress).toInt().coerceAtLeast(2)
    val animatedPoints = screenPoints.take(animatedPointsCount)

    if (animatedPoints.size >= 2) {
        // Draw line segments
        for (i in 0 until animatedPoints.size - 1) {
            drawLine(
                color = color,
                start = animatedPoints[i],
                end = animatedPoints[i + 1],
                strokeWidth = 3.dp.toPx(),
                cap = StrokeCap.Round
            )
        }

        // Draw data points
        animatedPoints.forEach { point ->
            drawCircle(
                color = color,
                radius = 6.dp.toPx(),
                center = point
            )
            drawCircle(
                color = Color.White,
                radius = 3.dp.toPx(),
                center = point
            )
        }
    }
}