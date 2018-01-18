package model

type DetailRequest struct {
	Id 			  string      `json:"id"`
	IdType	      string      `json:"idType"`
	Source        string      `json:"source"`
	Country		  string      `json:"country"`
}

func NewDetailRequest(id, idType, source, country string) *DetailRequest {
	o := &DetailRequest{
		Id: id,
		IdType: idType,
		Source: source,
		Country: country,
	}
	return o
}

// Checks if request is valid
func (r *DetailRequest) IsValid() bool {
	if r.Id == "" || r.IdType == "" || r.Source == "" {
		return false
	}

	return true
}
