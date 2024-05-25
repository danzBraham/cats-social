package match_entity

type Status string

const (
	Pending  Status = "pending"
	Approved Status = "approved"
	Rejected Status = "rejected"
)

type MatchCat struct {
	Id         string `json:"id"`
	MatchCatId string `json:"matchCatId"`
	UserCatId  string `json:"userCatId"`
	Message    string `json:"message"`
	Status     Status `json:"status"`
	IssuedBy   string `json:"issuedBy"`
	CreatedAt  string `json:"createdAt"`
}

type MatchCatRequest struct {
	MatchCatId string `json:"matchCatId" validate:"required,len=26"`
	UserCatId  string `json:"userCatId" validate:"required,len=26"`
	Message    string `json:"message" validate:"required,min=5,max=120"`
	Issuer     string `json:"userId" validate:"required,len=26"`
}
