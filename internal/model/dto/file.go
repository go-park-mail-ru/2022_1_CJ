package dto

type UploadFileResponse struct {
	URL string `json:"url"`
}

type GetFileRequest struct {
	URL string `json:"url"`
}
