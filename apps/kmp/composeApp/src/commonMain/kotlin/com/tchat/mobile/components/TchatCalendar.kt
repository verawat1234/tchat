package com.tchat.mobile.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.designsystem.TchatColors
import kotlinx.datetime.*
import kotlinx.datetime.format.*

/**
 * TchatCalendar - Date picker calendar component
 *
 * Features:
 * - Month/Year navigation with smooth transitions
 * - Date selection (single/range) with visual feedback
 * - Event indicators with custom styling
 * - Localized date formatting and week start
 */

data class CalendarEvent(
    val date: LocalDate,
    val title: String,
    val color: Color = TchatColors.primary,
    val priority: EventPriority = EventPriority.NORMAL
)

enum class EventPriority {
    LOW,
    NORMAL,
    HIGH,
    URGENT
}

enum class CalendarSelectionMode {
    SINGLE,
    RANGE,
    MULTIPLE
}

data class DateRange(
    val start: LocalDate?,
    val end: LocalDate?
) {
    val isValid: Boolean
        get() = start != null && end != null && start <= end
}

@Composable
fun TchatCalendar(
    selectedDate: LocalDate? = null,
    selectedRange: DateRange? = null,
    selectedDates: Set<LocalDate> = emptySet(),
    onDateSelected: (LocalDate) -> Unit = {},
    onRangeSelected: (DateRange) -> Unit = {},
    onDatesSelected: (Set<LocalDate>) -> Unit = {},
    selectionMode: CalendarSelectionMode = CalendarSelectionMode.SINGLE,
    events: List<CalendarEvent> = emptyList(),
    minDate: LocalDate? = null,
    maxDate: LocalDate? = null,
    modifier: Modifier = Modifier,
    showWeekNumbers: Boolean = false,
    weekStartsOnSunday: Boolean = false
) {
    var currentMonth by remember { mutableStateOf(Clock.System.todayIn(TimeZone.currentSystemDefault())) }
    val today = Clock.System.todayIn(TimeZone.currentSystemDefault())

    Column(
        modifier = modifier.fillMaxWidth()
    ) {
        // Calendar header with navigation
        CalendarHeader(
            currentMonth = currentMonth,
            onPreviousMonth = {
                currentMonth = currentMonth.minus(1, DateTimeUnit.MONTH)
            },
            onNextMonth = {
                currentMonth = currentMonth.plus(1, DateTimeUnit.MONTH)
            },
            onMonthYearClick = {
                // TODO: Implement month/year picker
            }
        )

        // Days of week header
        CalendarWeekHeader(weekStartsOnSunday = weekStartsOnSunday)

        // Calendar grid
        CalendarGrid(
            currentMonth = currentMonth,
            today = today,
            selectedDate = selectedDate,
            selectedRange = selectedRange,
            selectedDates = selectedDates,
            selectionMode = selectionMode,
            onDateSelected = onDateSelected,
            onRangeSelected = onRangeSelected,
            onDatesSelected = onDatesSelected,
            events = events,
            minDate = minDate,
            maxDate = maxDate,
            weekStartsOnSunday = weekStartsOnSunday,
            showWeekNumbers = showWeekNumbers
        )

        // Event list for selected date
        if (selectionMode == CalendarSelectionMode.SINGLE && selectedDate != null) {
            val dayEvents = events.filter { it.date == selectedDate }
            if (dayEvents.isNotEmpty()) {
                CalendarEventsList(
                    date = selectedDate,
                    events = dayEvents
                )
            }
        }
    }
}

@Composable
private fun CalendarHeader(
    currentMonth: LocalDate,
    onPreviousMonth: () -> Unit,
    onNextMonth: () -> Unit,
    onMonthYearClick: () -> Unit
) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        color = TchatColors.surface,
        shadowElevation = 1.dp
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            IconButton(
                onClick = onPreviousMonth,
                modifier = Modifier.semantics {
                    contentDescription = "Previous month"
                }
            ) {
                Icon(
                    imageVector = Icons.Default.ChevronLeft,
                    contentDescription = null
                )
            }

            TextButton(
                onClick = onMonthYearClick,
                modifier = Modifier.semantics {
                    contentDescription = "Select month and year"
                }
            ) {
                Text(
                    text = currentMonth.format(LocalDate.Format {
                        monthName(MonthNames.ENGLISH_FULL)
                        char(' ')
                        year()
                    }),
                    style = MaterialTheme.typography.headlineSmall,
                    fontWeight = FontWeight.Medium
                )
            }

            IconButton(
                onClick = onNextMonth,
                modifier = Modifier.semantics {
                    contentDescription = "Next month"
                }
            ) {
                Icon(
                    imageVector = Icons.Default.ChevronRight,
                    contentDescription = null
                )
            }
        }
    }
}

@Composable
private fun CalendarWeekHeader(weekStartsOnSunday: Boolean) {
    val daysOfWeek = if (weekStartsOnSunday) {
        listOf("Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat")
    } else {
        listOf("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun")
    }

    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 8.dp, vertical = 8.dp)
    ) {
        daysOfWeek.forEach { day ->
            Box(
                modifier = Modifier.weight(1f),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = day,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant,
                    fontWeight = FontWeight.Medium,
                    textAlign = TextAlign.Center
                )
            }
        }
    }
}

@Composable
private fun CalendarGrid(
    currentMonth: LocalDate,
    today: LocalDate,
    selectedDate: LocalDate?,
    selectedRange: DateRange?,
    selectedDates: Set<LocalDate>,
    selectionMode: CalendarSelectionMode,
    onDateSelected: (LocalDate) -> Unit,
    onRangeSelected: (DateRange) -> Unit,
    onDatesSelected: (Set<LocalDate>) -> Unit,
    events: List<CalendarEvent>,
    minDate: LocalDate?,
    maxDate: LocalDate?,
    weekStartsOnSunday: Boolean,
    showWeekNumbers: Boolean
) {
    val firstDayOfMonth = LocalDate(currentMonth.year, currentMonth.month, 1)
    val lastDayOfMonth = currentMonth.atEndOfMonth()
    val firstDayOfWeek = if (weekStartsOnSunday) DayOfWeek.SUNDAY else DayOfWeek.MONDAY

    // Calculate calendar days
    val calendarDays = buildList {
        // Add previous month days
        val daysFromPrevMonth = ((firstDayOfMonth.dayOfWeek.ordinal - firstDayOfWeek.ordinal + 7) % 7)
        for (i in daysFromPrevMonth downTo 1) {
            add(firstDayOfMonth.minus(i, DateTimeUnit.DAY))
        }

        // Add current month days
        var currentDay = firstDayOfMonth
        while (currentDay <= lastDayOfMonth) {
            add(currentDay)
            currentDay = currentDay.plus(1, DateTimeUnit.DAY)
        }

        // Add next month days to fill the grid
        val totalDays = size
        val remainingDays = 42 - totalDays // 6 rows * 7 days
        var nextMonthDay = lastDayOfMonth.plus(1, DateTimeUnit.DAY)
        repeat(remainingDays) {
            add(nextMonthDay)
            nextMonthDay = nextMonthDay.plus(1, DateTimeUnit.DAY)
        }
    }

    LazyVerticalGrid(
        columns = GridCells.Fixed(7),
        modifier = Modifier.fillMaxWidth(),
        contentPadding = PaddingValues(horizontal = 8.dp, vertical = 4.dp),
        verticalArrangement = Arrangement.spacedBy(2.dp),
        horizontalArrangement = Arrangement.spacedBy(2.dp)
    ) {
        items(calendarDays) { date ->
            CalendarDay(
                date = date,
                isCurrentMonth = date.month == currentMonth.month,
                isToday = date == today,
                isSelected = when (selectionMode) {
                    CalendarSelectionMode.SINGLE -> date == selectedDate
                    CalendarSelectionMode.RANGE -> selectedRange?.let { range ->
                        date == range.start || date == range.end ||
                        (range.isValid && date > range.start!! && date < range.end!!)
                    } ?: false
                    CalendarSelectionMode.MULTIPLE -> date in selectedDates
                },
                isEnabled = (minDate == null || date >= minDate) &&
                           (maxDate == null || date <= maxDate),
                events = events.filter { it.date == date },
                onClick = {
                    when (selectionMode) {
                        CalendarSelectionMode.SINGLE -> onDateSelected(date)
                        CalendarSelectionMode.RANGE -> {
                            val currentRange = selectedRange
                            when {
                                currentRange?.start == null -> {
                                    onRangeSelected(DateRange(date, null))
                                }
                                currentRange.end == null -> {
                                    val newRange = if (date >= currentRange.start!!) {
                                        DateRange(currentRange.start, date)
                                    } else {
                                        DateRange(date, currentRange.start)
                                    }
                                    onRangeSelected(newRange)
                                }
                                else -> {
                                    onRangeSelected(DateRange(date, null))
                                }
                            }
                        }
                        CalendarSelectionMode.MULTIPLE -> {
                            val newSelection = if (date in selectedDates) {
                                selectedDates - date
                            } else {
                                selectedDates + date
                            }
                            onDatesSelected(newSelection)
                        }
                    }
                }
            )
        }
    }
}

@Composable
private fun CalendarDay(
    date: LocalDate,
    isCurrentMonth: Boolean,
    isToday: Boolean,
    isSelected: Boolean,
    isEnabled: Boolean,
    events: List<CalendarEvent>,
    onClick: () -> Unit
) {
    val backgroundColor = when {
        isSelected -> TchatColors.primary
        isToday -> TchatColors.primary.copy(alpha = 0.1f)
        else -> Color.Transparent
    }

    val textColor = when {
        !isEnabled -> TchatColors.disabled
        isSelected -> TchatColors.onPrimary
        !isCurrentMonth -> TchatColors.onSurfaceVariant.copy(alpha = 0.6f)
        isToday -> TchatColors.primary
        else -> TchatColors.onSurface
    }

    Box(
        modifier = Modifier
            .aspectRatio(1f)
            .clip(RoundedCornerShape(8.dp))
            .background(backgroundColor)
            .clickable(enabled = isEnabled) { onClick() }
            .semantics {
                contentDescription = buildString {
                    append(date.format(LocalDate.Format {
                        monthName(MonthNames.ENGLISH_FULL)
                        char(' ')
                        dayOfMonth()
                        char(',')
                        char(' ')
                        year()
                    }))
                    if (isToday) append(", today")
                    if (isSelected) append(", selected")
                    if (events.isNotEmpty()) {
                        append(", ${events.size} event${if (events.size > 1) "s" else ""}")
                    }
                }
                role = Role.Button
                if (!isEnabled) disabled()
            }
            .padding(4.dp),
        contentAlignment = Alignment.Center
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(2.dp)
        ) {
            Text(
                text = date.dayOfMonth.toString(),
                style = MaterialTheme.typography.bodyMedium,
                color = textColor,
                fontWeight = if (isToday || isSelected) FontWeight.Medium else FontWeight.Normal,
                fontSize = 14.sp
            )

            // Event indicators
            if (events.isNotEmpty() && events.size <= 3) {
                Row(
                    horizontalArrangement = Arrangement.spacedBy(1.dp)
                ) {
                    events.take(3).forEach { event ->
                        Box(
                            modifier = Modifier
                                .size(4.dp)
                                .clip(CircleShape)
                                .background(
                                    if (isSelected) TchatColors.onPrimary.copy(alpha = 0.8f)
                                    else event.color
                                )
                        )
                    }
                }
            } else if (events.size > 3) {
                Text(
                    text = "${events.size}",
                    style = MaterialTheme.typography.bodySmall,
                    color = if (isSelected) TchatColors.onPrimary else TchatColors.primary,
                    fontSize = 10.sp
                )
            }
        }
    }
}

@Composable
private fun CalendarEventsList(
    date: LocalDate,
    events: List<CalendarEvent>
) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        Text(
            text = "Events for ${date.format(LocalDate.Format {
                monthName(MonthNames.ENGLISH_FULL)
                char(' ')
                dayOfMonth()
            })}",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Medium
        )

        events.forEach { event ->
            Surface(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(8.dp),
                color = event.color.copy(alpha = 0.1f),
                border = androidx.compose.foundation.BorderStroke(
                    1.dp,
                    event.color.copy(alpha = 0.3f)
                )
            ) {
                Row(
                    modifier = Modifier.padding(12.dp),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Box(
                        modifier = Modifier
                            .size(8.dp)
                            .clip(CircleShape)
                            .background(event.color)
                    )

                    Text(
                        text = event.title,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurface,
                        modifier = Modifier.weight(1f)
                    )

                    if (event.priority != EventPriority.NORMAL) {
                        Icon(
                            imageVector = when (event.priority) {
                                EventPriority.HIGH -> Icons.Default.PriorityHigh
                                EventPriority.URGENT -> Icons.Default.NotificationImportant
                                EventPriority.LOW -> Icons.Default.ArrowDownward
                                else -> Icons.Default.Circle
                            },
                            contentDescription = "${event.priority} priority",
                            tint = event.color,
                            modifier = Modifier.size(16.dp)
                        )
                    }
                }
            }
        }
    }
}

// Extension function to get the last day of month
private fun LocalDate.atEndOfMonth(): LocalDate {
    val daysInMonth = when (month) {
        Month.FEBRUARY -> if (isLeapYear(year)) 29 else 28
        Month.APRIL, Month.JUNE, Month.SEPTEMBER, Month.NOVEMBER -> 30
        else -> 31
    }
    return LocalDate(year, month, daysInMonth)
}

private fun isLeapYear(year: Int): Boolean {
    return year % 4 == 0 && (year % 100 != 0 || year % 400 == 0)
}