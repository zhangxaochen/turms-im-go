package model

import (
	"bytes"
	"encoding/json"
)

type BulkRequest struct {
	Operations []interface{}
}

// Serialize serializes the bulk request into an NDJSON format.
// @MappedFrom serialize(BulkRequest value, JsonGenerator gen, SerializerProvider serializers)
func (r *BulkRequest) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	for _, op := range r.Operations {
		data, err := json.Marshal(op)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}

type BulkResponse struct {
	Errors bool               `json:"errors"`
	Items  []BulkResponseItem `json:"items"`
}

type BulkResponseRecord struct {
	Index *BulkResponseItem `json:"index"`
}

type BulkResponseItem struct {
	ID      *string          `json:"_id"`
	Index   string           `json:"_index"`
	Status  int              `json:"status"`
	Error   *ErrorCause      `json:"error,omitempty"`
	Result  *string          `json:"result,omitempty"`
	SeqNo   *int64           `json:"_seq_no,omitempty"`
	Shards  *ShardStatistics `json:"_shards,omitempty"`
	Version *int64           `json:"_version,omitempty"`
}

type ClosePointInTimeRequest struct {
	ID string `json:"id"`
}

type CreateIndexRequest struct {
	Mappings *TypeMapping     `json:"mappings,omitempty"`
	Settings *json.RawMessage `json:"settings,omitempty"`
}

type DeleteByQueryRequest struct {
	Query json.RawMessage `json:"query"`
}

type DeleteByQueryResponse struct {
	Deleted bool `json:"deleted"`
}

type DeleteResponse struct {
	Result string `json:"result"`
}

type ErrorCause struct {
	Type       *string      `json:"type,omitempty"`
	Reason     *string      `json:"reason,omitempty"`
	StackTrace *string      `json:"stack_trace,omitempty"`
	CausedBy   *ErrorCause  `json:"caused_by,omitempty"`
	RootCause  []ErrorCause `json:"root_cause,omitempty"`
	Suppressed []ErrorCause `json:"suppressed,omitempty"`
}

type ErrorResponse struct {
	Error ErrorCause `json:"error"`
}

type FieldCollapse struct {
	Field string `json:"field"`
}

type HealthResponse struct {
	ClusterName string `json:"cluster_name"`
	Status      string `json:"status"`
}

type Highlight struct {
	Fields map[string]interface{} `json:"fields"`
}

type IndexSettings struct {
	Analysis *IndexSettingsAnalysis `json:"analysis,omitempty"`
}

type IndexSettingsAnalysis struct {
	Analyzer   map[string]map[string]interface{} `json:"analyzer,omitempty"`
	CharFilter map[string]map[string]interface{} `json:"char_filter,omitempty"`
	Filter     map[string]map[string]interface{} `json:"filter,omitempty"`
	Normalizer map[string]map[string]interface{} `json:"normalizer,omitempty"`
	Tokenizer  map[string]map[string]interface{} `json:"tokenizer,omitempty"`
}

// Merge merges another IndexSettingsAnalysis into this one.
// @MappedFrom merge(IndexSettingsAnalysis analysis)
func (a *IndexSettingsAnalysis) Merge(other *IndexSettingsAnalysis) {
	if other == nil {
		return
	}
	a.Analyzer = mergeMaps(a.Analyzer, other.Analyzer)
	a.CharFilter = mergeMaps(a.CharFilter, other.CharFilter)
	a.Filter = mergeMaps(a.Filter, other.Filter)
	a.Normalizer = mergeMaps(a.Normalizer, other.Normalizer)
	a.Tokenizer = mergeMaps(a.Tokenizer, other.Tokenizer)
}

func mergeMaps(m1, m2 map[string]map[string]interface{}) map[string]map[string]interface{} {
	if m1 == nil {
		return m2
	}
	if m2 == nil {
		return m1
	}
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

type PointInTimeReference struct {
	ID        string `json:"id"`
	KeepAlive string `json:"keep_alive"`
}

type Property struct {
	Type           PropertyType        `json:"type"`
	Analyzer       *string             `json:"analyzer,omitempty"`
	SearchAnalyzer *string             `json:"search_analyzer,omitempty"`
	Fields         map[string]Property `json:"fields,omitempty"`
}

type PropertyType string

const (
	PropertyTypeKeyword PropertyType = "keyword"
	PropertyTypeText    PropertyType = "text"
)

type Script struct {
	Source string                 `json:"source"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type SearchRequest struct {
	From           *int                  `json:"from,omitempty"`
	Size           *int                  `json:"size,omitempty"`
	Query          json.RawMessage       `json:"query"`
	Highlight      *Highlight            `json:"highlight,omitempty"`
	Collapse       *FieldCollapse        `json:"collapse,omitempty"`
	PointInTime    *PointInTimeReference `json:"pit,omitempty"`
	Sort           []interface{}         `json:"sort,omitempty"`
	TrackTotalHits bool                  `json:"track_total_hits,omitempty"`
}

type ShardFailure struct {
	Index  *string    `json:"index,omitempty"`
	Node   *string    `json:"node,omitempty"`
	Reason ErrorCause `json:"reason"`
	Shard  int        `json:"shard"`
	Status *string    `json:"status,omitempty"`
}

type ShardStatistics struct {
	Failed     int64          `json:"failed"`
	Successful int64          `json:"successful"`
	Total      int64          `json:"total"`
	Failures   []ShardFailure `json:"failures,omitempty"`
	Skipped    *int64         `json:"skipped,omitempty"`
}

type TypeMapping struct {
	Properties map[string]Property `json:"properties"`
}

type UpdateByQueryRequest struct {
	Query json.RawMessage `json:"query"`
}

type UpdateByQueryResponse struct {
	Updated int64 `json:"updated"`
}
