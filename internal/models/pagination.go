package models

type Pagination struct {
	Page         int
	PerPage      int
	TotalPages   int
	TotalCount   int
	SelfLink     string
	FirstLink    string
	PreviousLink string
	NextLink     string
	LastLink     string
}
