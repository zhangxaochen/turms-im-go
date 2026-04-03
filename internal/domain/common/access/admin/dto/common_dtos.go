package dto

import "time"

// @MappedFrom StatisticsRecordDTO(Date date, Long total)
type StatisticsRecordDTO struct {
	Date  time.Time `json:"date"`
	Total int64     `json:"total"`
}
