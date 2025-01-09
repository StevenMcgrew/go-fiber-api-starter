package models

type Pagination struct {
	Page         uint   `json:"page"`
	PerPage      uint   `json:"perPage"`
	TotalPages   uint   `json:"totalPages"`
	TotalCount   uint   `json:"totalCount"`
	SelfLink     string `json:"selfLink"`
	FirstLink    string `json:"firstLink"`
	PreviousLink string `json:"previousLink"`
	NextLink     string `json:"nextLink"`
	LastLink     string `json:"lastLink"`
}
