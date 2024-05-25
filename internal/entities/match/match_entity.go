package match_entity

import cat_entity "github.com/danzbraham/cats-social/internal/entities/cat"

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

type IssuerDetail struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type GetMatchCatResponse struct {
	Id             string                   `json:"id"`
	IssuedBy       IssuerDetail             `json:"issuedBy"`
	MatchCatDetail cat_entity.GetCatReponse `json:"matchCatDetail"`
	UserCatDetail  cat_entity.GetCatReponse `json:"userCatDetail"`
	Message        string                   `json:"message"`
	CreatedAt      string                   `json:"createdAt"`
}
