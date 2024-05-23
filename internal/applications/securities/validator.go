package securities

type Validator interface {
	ValidatePayload(payload interface{}) error
}
