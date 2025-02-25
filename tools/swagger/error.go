package swagger

type ProtobufAny struct {
	Type string `json:"@type"`
}

type RpcStatus struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Details []ProtobufAny `json:"details"`
}
