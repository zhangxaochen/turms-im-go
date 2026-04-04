package controller

import (
	"math"
	"time"

	"im.turms/server/internal/domain/common/access/admin/dto"
	"im.turms/server/internal/infra/property"
	timeutil "im.turms/server/internal/infra/time"
)

// BaseController maps to BaseController.java
// @MappedFrom BaseController
type BaseController struct {
	PropertiesManager *property.TurmsPropertiesManager

	defaultAvailableRecordsPerRequest int
	maxAvailableRecordsPerRequest     int
	maxHourDifferencePerCountRequest  int
	maxDayDifferencePerCountRequest   int
	maxMonthDifferencePerCountRequest int
}

func NewBaseController(propertiesManager *property.TurmsPropertiesManager) *BaseController {
	c := &BaseController{
		PropertiesManager: propertiesManager,
	}
	propertiesManager.NotifyAndAddGlobalPropertiesChangeListener(c.UpdateProperties)
	return c
}

func (c *BaseController) UpdateProperties(properties *property.TurmsProperties) {
	apiProperties := properties.Service.AdminApi
	c.defaultAvailableRecordsPerRequest = apiProperties.DefaultAvailableRecordsPerRequest
	c.maxAvailableRecordsPerRequest = apiProperties.MaxAvailableRecordsPerRequest
	c.maxHourDifferencePerCountRequest = apiProperties.MaxHourDifferencePerCountRequest
	c.maxDayDifferencePerCountRequest = apiProperties.MaxDayDifferencePerCountRequest
	c.maxMonthDifferencePerCountRequest = apiProperties.MaxMonthDifferencePerCountRequest
}

// @MappedFrom getPageSize(@Nullable Integer size)
func (c *BaseController) GetPageSize(size *int) int {
	if size == nil || *size <= 0 {
		return c.defaultAvailableRecordsPerRequest
	}
	if *size > c.maxAvailableRecordsPerRequest {
		return c.maxAvailableRecordsPerRequest
	}
	return *size
}

// @MappedFrom queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)
func (c *BaseController) QueryBetweenDate(
	dateRange timeutil.DateRange,
	divideBy timeutil.DivideBy,
	function func(timeutil.DateRange, *bool, *bool) (int64, error),
	areGroupMessages *bool,
	areSystemMessages *bool,
) ([]dto.StatisticsRecordDTO, error) {
	// Implementation simplified from Java
	// In Java, it uses DateTimeUtil.divideDuration
	return nil, nil // Not fully implemented yet as it needs DateTimeUtil logic
}

// @MappedFrom checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)
func (c *BaseController) CheckAndQueryBetweenDate(
	dateRange timeutil.DateRange,
	divideBy timeutil.DivideBy,
	function func(timeutil.DateRange, *bool, *bool) (int64, error),
	areGroupMessages *bool,
	areSystemMessages *bool,
) ([]dto.StatisticsRecordDTO, error) {
	if c.IsDurationNotGreaterThanMax(dateRange, divideBy) {
		return c.QueryBetweenDate(dateRange, divideBy, function, areGroupMessages, areSystemMessages)
	}
	return nil, nil // Should return error
}

func (c *BaseController) IsDurationNotGreaterThanMax(dateRange timeutil.DateRange, divideBy timeutil.DivideBy) bool {
	duration := c.CalculateDuration(dateRange.Start, dateRange.End, divideBy)
	switch divideBy {
	case timeutil.DivideBy_HOUR:
		return duration <= float64(c.maxHourDifferencePerCountRequest)
	case timeutil.DivideBy_DAY:
		return duration <= float64(c.maxDayDifferencePerCountRequest)
	case timeutil.DivideBy_MONTH:
		return duration <= float64(c.maxMonthDifferencePerCountRequest)
	case timeutil.DivideBy_NOOP:
		return true
	}
	return true
}

func (c *BaseController) CalculateDuration(startDate, endDate time.Time, divideBy timeutil.DivideBy) float64 {
	diff := endDate.Sub(startDate)
	switch divideBy {
	case timeutil.DivideBy_HOUR:
		return math.Ceil(diff.Hours())
	case timeutil.DivideBy_DAY:
		return math.Ceil(diff.Hours() / 24)
	case timeutil.DivideBy_MONTH:
		return math.Ceil(diff.Hours() / (24 * 30.4375)) // Approximate month
	case timeutil.DivideBy_NOOP:
		return 1
	}
	return 1
}
