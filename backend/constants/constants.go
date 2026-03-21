package constants

// Constants
const (
	// Environment modes
	EnvModeDev  = "dev"
	EnvModeProd = "prod"

	// Collection name suffixes - will be prefixed with product name
	CollectionSuffixResources = "_resources"
	CollectionSuffixTags      = "_tags"

	// Valid product names
	ProductEcomm = "ecomm"

	DefaultPageSize = 20
	MaxPageSize     = 100
	MaxFileSize     = 500 << 20 // 500MB

	ProductContextKey = "product"
	ProductParamKey   = "product"

	// Resource Types
	ResourceTypeVideo   = "video"
	ResourceTypePDF     = "pdf"
	ResourceTypeArticle = "article"
	ResourceTypeImage   = "image"

	// Query Parameter Names
	QueryParamType   = "type"
	QueryParamTags   = "tags"
	QueryParamSearch = "search"
	QueryParamCursor = "cursor"
	QueryParamLimit  = "limit"

	// Form Field Names
	FormFieldTitle        = "title"
	FormFieldDescription  = "description"
	FormFieldType         = "type"
	FormFieldURL          = "url"
	FormFieldThumbnailURL = "thumbnailUrl"
	FormFieldTags         = "tags"
	FormFieldFile         = "file"
	FormFieldThumbnail    = "thumbnail"

	// Default Values
	DefaultLimitValue = "20"

	// URL expiration times
	DefaultSignedURLExpiration = 60 // 1 hour

	// Error message prefixes
	ErrFileValidationFailed = "file validation failed"
)

// ResourceTypes ...
var ResourceTypes = []string{
	ResourceTypeVideo,
	ResourceTypePDF,
	ResourceTypeArticle,
}

// ValidProducts ...
var ValidProducts = []string{
	ProductEcomm,
}

// GetResourcesCollectionName returns the collection name for resources for a given productMore actions
// product_name + "_resources"
func GetResourcesCollectionName(product string) string {
	return product + CollectionSuffixResources
}

// GetTagsCollectionName returns the collection name for tags for a given productMore actions
// product_name + "_tags"
func GetTagsCollectionName(product string) string {
	return product + CollectionSuffixTags
}
