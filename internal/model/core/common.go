package core

type PaginationParameters struct {
	Limit int64 `query:"limit,omitempty"`
	Page  int64 `query:"page,omitempty"`
}
