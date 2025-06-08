package constants

// Constants
const (
	CollectionResources = "resources"
	CollectionTags      = "tags"
	DefaultPageSize     = 20
	MaxPageSize         = 100
	MaxFileSize         = 100 << 20 // 100MB

	AdminSecretHeader        = "AdminSecret"
	AdminSecretQueryParamKey = "adminSecret"

	// Resource Types
	ResourceTypeVideo   = "video"
	ResourceTypePDF     = "pdf"
	ResourceTypeArticle = "article"

	// Generic errors
	InvalidParam       = "invalid_param"
	InvalidContentType = "invalid_content_type"
	InvalidPayload     = "invalid_payload"
	Unauthorized       = "unauthorized"

	// DB specific errors
	QueryFailed          = "query_failed"
	MutationFailed       = "mutation_failed"
	NotFound             = "not_found"
	DataConversionFailed = "data_conversion_failed"
	UploadFailed         = "upload_failed"
)

// ResourceTypes ...
var ResourceTypes = []string{
	ResourceTypeVideo,
	ResourceTypePDF,
	ResourceTypeArticle,
}
