package aisera

type QueryOutput struct {
	Data       interface{} `json:"data"`
	TotalCount int         `json:"totalCount"`
}

type ExecutionParam struct {
	ID           int    `json:"id"`
	EntityTypeID string `json:"entityTypeId"`
	Value        string `json:"value"`
}

type User struct {
	Name      any    `json:"name"`
	UserEmail string `json:"user_email"`
}

type Execution struct {
	PipelineStatus string          `json:"pipeline_status"`
	Output         any             `json:"output"`
	Message        any             `json:"message"`
	Error          any             `json:"error"`
	StartedAt      string          `json:"started_at"`
	Duration       any             `json:"_duration"`
	ID             int             `json:"id"`
	EnvVars        GenericKeyValue `json:"env_vars,omitempty"`
	Job            struct {
		ReferenceID string `json:"reference_id"`
	} `json:"job"`
	Params []ExecutionParam `json:"_params"`
	User   User             `json:"_user,omitempty"`
}
