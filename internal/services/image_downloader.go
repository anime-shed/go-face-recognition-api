package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"strings"
	"time"

	// Import image decoders
	_ "image/jpeg"
	_ "image/png"

	"github.com/sirupsen/logrus"

	"face-recognition-api/internal/config"
	"face-recognition-api/internal/models"
)

// ImageDownloader handles downloading and validating images from URLs
type ImageDownloader struct {
	client *http.Client
	config config.LimitsConfig
	logger *logrus.Logger
}

// NewImageDownloader creates a new image downloader instance
func NewImageDownloader(cfg config.LimitsConfig, logger *logrus.Logger) *ImageDownloader {
	return &ImageDownloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Internal service - skip certificate validation
				},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  false,
			},
		},
		config: cfg,
		logger: logger,
	}
}

// DownloadImage downloads an image from the given URL and returns the decoded image
func (id *ImageDownloader) DownloadImage(ctx context.Context, imageURL string) (image.Image, models.ImageMetadata, error) {
	// Validate URL format
	if err := id.validateURL(imageURL); err != nil {
		return nil, models.ImageMetadata{}, fmt.Errorf("invalid URL: %w", err)
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, models.ImageMetadata{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Face-Recognition-API/1.0")
	req.Header.Set("Accept", "image/jpeg,image/png,image/*")

	// Execute request
	resp, err := id.client.Do(req)
	if err != nil {
		return nil, models.ImageMetadata{}, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, models.ImageMetadata{}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !id.isValidImageType(contentType) {
		return nil, models.ImageMetadata{}, fmt.Errorf("unsupported content type: %s", contentType)
	}

	// Check content length
	if resp.ContentLength > id.config.MaxImageSize {
		return nil, models.ImageMetadata{}, fmt.Errorf("image too large: %d bytes (max: %d)", resp.ContentLength, id.config.MaxImageSize)
	}

	// Decode image
	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, models.ImageMetadata{}, fmt.Errorf("failed to decode image: %w", err)
	}

	// Validate image dimensions
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	
	if width > id.config.MaxWidth || height > id.config.MaxHeight {
		return nil, models.ImageMetadata{}, fmt.Errorf("image dimensions too large: %dx%d (max: %dx%d)", 
			width, height, id.config.MaxWidth, id.config.MaxHeight)
	}

	// Create metadata
	metadata := models.ImageMetadata{
		Width:     width,
		Height:    height,
		Format:    strings.ToUpper(format),
		SizeBytes: resp.ContentLength,
		URL:       imageURL,
	}

	id.logger.WithFields(logrus.Fields{
		"url":         imageURL,
		"width":       width,
		"height":      height,
		"format":      format,
		"size_bytes":  resp.ContentLength,
	}).Info("Image downloaded successfully")

	return img, metadata, nil
}

// validateURL validates the image URL format and security
func (id *ImageDownloader) validateURL(imageURL string) error {
	if imageURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	// Check host
	if parsedURL.Host == "" {
		return fmt.Errorf("missing host in URL")
	}

	// Basic SSRF protection - block private IP ranges
	if id.isPrivateIP(parsedURL.Host) {
		return fmt.Errorf("access to private IP ranges is not allowed")
	}

	return nil
}

// isValidImageType checks if the content type is a supported image format
func (id *ImageDownloader) isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}

	return false
}

// isPrivateIP performs basic check for private IP ranges (simplified)
func (id *ImageDownloader) isPrivateIP(host string) bool {
	// This is a simplified check - in production, you'd want more comprehensive validation
	privateHosts := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
	}

	host = strings.ToLower(host)
	for _, privateHost := range privateHosts {
		if strings.Contains(host, privateHost) {
			return true
		}
	}

	// Check for private IP prefixes
	if strings.HasPrefix(host, "10.") ||
		strings.HasPrefix(host, "192.168.") ||
		strings.HasPrefix(host, "172.") {
		return true
	}

	return false
}