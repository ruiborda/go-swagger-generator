package openapi_spec

type PathItemEntity struct {
	Get        *OperationEntity  `json:"get,omitempty"`
	Post       *OperationEntity  `json:"post,omitempty"`
	Put        *OperationEntity  `json:"put,omitempty"`
	Delete     *OperationEntity  `json:"delete,omitempty"`
	Options    *OperationEntity  `json:"options,omitempty"`
	Head       *OperationEntity  `json:"head,omitempty"`
	Patch      *OperationEntity  `json:"patch,omitempty"`
	Parameters []ParameterEntity `json:"parameters,omitempty"`
	Ref        string            `json:"$ref,omitempty"`
}
