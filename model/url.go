package model

type UrlStatus string

const (
	StatusNew        UrlStatus = "NEW"
	StatusProcessing UrlStatus = "PROCESSING"
	StatusDone       UrlStatus = "DONE"
	StatusError      UrlStatus = "ERROR"
)

type URL struct {
	Id       int
	Url      string
	Status   UrlStatus
	HttpCode int32
}
