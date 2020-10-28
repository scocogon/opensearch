package response

// Response 应答
type Response struct {
	Status         string      `json:"status"`
	RequestID      string      `json:"request_id"`
	Result         Result      `json:"result"`
	Errors         []ErrorNode `json:"errors"`
	Tracer         string      `json:"tracer"`
	OpsRequestMisc string      `josn:"ops_request_misc"`
}

func (r *Response) Success() bool {
	return r.Status == "OK"
}

// Result response.result
type Result struct {
	SearchTime  float64                  `json:"searchtime"`
	Total       int                      `json:"total"`
	Num         int                      `json:"num"`
	ViewTotal   int                      `json:"viewtotal"`
	ComputeCost []ComputeCost            `json:"compute_cost"`
	Items       []map[string]interface{} `json:"items"`
}

// ComputeCost result.compute_cost
type ComputeCost struct {
	IndexName string  `json:"index_name"`
	Value     float64 `json:"value"`
}

// ErrorNode errors
type ErrorNode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
