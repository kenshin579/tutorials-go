package opensearch

import "time"

// Product는 검색 쿼리 학습용 상품 모델이다.
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	InStock     bool      `json:"in_stock"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
}

// AccessLog는 Dashboards 시각화용 API 로그 모델이다.
type AccessLog struct {
	Timestamp      time.Time `json:"timestamp"`
	Method         string    `json:"method"`
	Endpoint       string    `json:"endpoint"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs float64   `json:"response_time_ms"`
	ErrorMessage   string    `json:"error_message"`
	ClientIP       string    `json:"client_ip"`
	UserAgent      string    `json:"user_agent"`
	RequestBody    string    `json:"request_body"`
	ServiceName    string    `json:"service_name"`
}
