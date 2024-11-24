package http

import (
	"errors"
	"strconv"
	"time"

	"github.com/Coke3a/HotelManagement/internal/core/domain"
	"github.com/Coke3a/HotelManagement/internal/core/port"
	"github.com/gin-gonic/gin"
)

// LogHandler represents the HTTP handler for log-related requests
type LogHandler struct {
	svc port.LogService
}

// NewLogHandler creates a new LogHandler instance
func NewLogHandler(svc port.LogService) *LogHandler {
	return &LogHandler{
		svc,
	}
}

type listLogsRequest struct {
	Skip uint64 `query:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `query:"limit" binding:"required,min=1" example:"10"`
}

func (lh *LogHandler) GetLogs(ctx *gin.Context) {
	var req listLogsRequest
	var logsList []logResponse

	skip := ctx.Query("skip")
	limit := ctx.Query("limit")

	skipUint, err := strconv.ParseUint(skip, 10, 64)
	if err != nil {
		validationError(ctx, err)
		return
	}

	limitUint, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		validationError(ctx, err)
		return
	}

	logs, totalCount, err := lh.svc.GetLogs(ctx, skipUint, limitUint)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, log := range logs {
		rsp, err := newLogResponse(&log)
		if err != nil {
			handleError(ctx, err)
			return
		}
		logsList = append(logsList, rsp)
	}

	meta := newMeta(totalCount, req.Limit, req.Skip)
	rsp := toMap(meta, logsList, "logs")

	handleSuccess(ctx, rsp)
}



// logResponse represents the response body for a log
type logResponse struct {
	ID        uint64 `json:"id"`
	TableName string `json:"table_name"`
	RecordID  uint64 `json:"record_id"`
	Action    string `json:"action"`
	UserID    uint64 `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// newLogResponse creates a new log response
func newLogResponse(log *domain.Log) (logResponse, error) {
	if log == nil {
		return logResponse{}, errors.New("log is nil")
	}

	createdAt := log.CreatedAt

	return logResponse{
		ID:        log.ID,
		TableName: log.TableName,
		RecordID:  log.RecordID,
		Action:    log.Action,
		UserID:    log.UserID,
		CreatedAt: createdAt,
	}, nil
}	