package http

import (
    "time"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/Coke3a/HotelManagement/internal/core/domain"
    "github.com/Coke3a/HotelManagement/internal/core/port"
)

type DailyBookingSummaryHandler struct {
    svc port.DailyBookingSummaryService
}

func NewDailyBookingSummaryHandler(svc port.DailyBookingSummaryService) *DailyBookingSummaryHandler {
    return &DailyBookingSummaryHandler{
        svc,
    }
}

// GenerateDailySummary godoc
// @Summary Generate daily booking summary
// @Description Generate summary for specified date
// @Tags daily-booking-summary
// @Accept json
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} domain.DailyBookingSummary
// @Router /api/v1/daily-summary/generate [post]
func (h *DailyBookingSummaryHandler) GenerateDailySummary(ctx *gin.Context) {
    dateStr := ctx.Query("date")
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        validationError(ctx, err)
        return
    }

    summary, err := h.svc.GenerateDailySummary(ctx, date)
    if err != nil {
        handleError(ctx, err)
        return
    }

    ctx.JSON(200, summary)
}

// UpdateSummaryStatus godoc
// @Summary Update daily booking summary status
// @Description Update status of a daily booking summary
// @Tags daily-booking-summary
// @Accept json
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Param status query int true "Status (0: Unchecked, 1: Checked, 2: Confirmed)"
// @Success 200 {object} domain.DailyBookingSummary
// @Router /api/v1/daily-summary/status [put]
func (h *DailyBookingSummaryHandler) UpdateSummaryStatus(ctx *gin.Context) {
    dateStr := ctx.Query("date")
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        validationError(ctx, err)
        return
    }

    statusStr := ctx.Query("status")
    status, err := strconv.Atoi(statusStr)
    if err != nil {
        validationError(ctx, err)
        return
    }

    summary, err := h.svc.UpdateSummaryStatus(ctx, date, domain.SummaryStatus(status))
    if err != nil {
        handleError(ctx, err)
        return
    }

    ctx.JSON(200, summary)
}

// GetSummaryByDate godoc
// @Summary Get daily booking summary by date
// @Description Get summary for specified date
// @Tags daily-booking-summary
// @Accept json
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} domain.DailyBookingSummary
// @Router /api/v1/daily-summary [get]
func (h *DailyBookingSummaryHandler) GetSummaryByDate(ctx *gin.Context) {
    dateStr := ctx.Query("date")
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        validationError(ctx, err)
        return
    }

    summary, err := h.svc.GetSummaryByDate(ctx, date)
    if err != nil {
        handleError(ctx, err)
        return
    }

    ctx.JSON(200, summary)
}

// ListSummaries godoc
// @Summary List daily booking summaries
// @Description Get paginated list of daily booking summaries
// @Tags daily-booking-summary
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} domain.DailyBookingSummary
// @Router /api/v1/daily-summary/list [get]
func (h *DailyBookingSummaryHandler) ListSummaries(ctx *gin.Context) {
    page, _ := strconv.ParseUint(ctx.DefaultQuery("page", "1"), 10, 64)
    limit, _ := strconv.ParseUint(ctx.DefaultQuery("limit", "10"), 10, 64)
    skip := (page - 1) * limit

    summaries, total, err := h.svc.ListSummaries(ctx, skip, limit)
    if err != nil {
        handleError(ctx, err)
        return
    }

    ctx.JSON(200, gin.H{
        "data":  summaries,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}