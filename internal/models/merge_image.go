package models

type MergeImageInput struct {
	BaseObName  string `json:"base_object" validate:"required"`
	FrameObject string `json:"frame_object" validate:"required"`
	Bucket      string `json:"bucket" validate:"required"`
	// UpdateBy     string `json:"update_by"`
}

// Kết quả chi tiết của việc merge ảnh (ví dụ link ảnh đã upload)
type MergeImageResult struct {
	URL string `json:"url"`
}
