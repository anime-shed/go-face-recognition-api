package models

// FaceDetectionRequest represents the request for face detection endpoint
type FaceDetectionRequest struct {
	ImageURL string `json:"image_url" binding:"required,url"`
}

// SelfieValidationRequest represents the request for selfie validation endpoint
type SelfieValidationRequest struct {
	ImageURL string `json:"image_url" binding:"required,url"`
	MinFaces int    `json:"min_faces" default:"1"`
	MaxFaces int    `json:"max_faces" default:"1"`
}

// VisualDetectionRequest represents the request for visual detection endpoint
type VisualDetectionRequest struct {
	ImageURL    string `json:"image_url" binding:"required,url"`
	CircleColor string `json:"circle_color" default:"red"`
	LineWidth   int    `json:"line_width" default:"3"`
}