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
	ProductAdmin = "admin"
	ProductCrm   = "crm"

	DefaultPageSize = 20
	MaxPageSize     = 100
	MaxFileSize     = 100 << 20 // 100MB

	AdminSecretHeader        = "AdminSecret"
	AdminSecretQueryParamKey = "adminSecret"

	ProductContextKey = "product"
	ProductParamKey   = "product"

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

// ValidProducts ...
var ValidProducts = []string{
	ProductEcomm,
	ProductAdmin,
	ProductCrm,
}

// GetResourcesCollectionName returns the collection name for resources for a given product
// product_name + "_resources"
func GetResourcesCollectionName(product string) string {
	return product + CollectionSuffixResources
}

// GetTagsCollectionName returns the collection name for tags for a given product
// product_name + "_tags"
func GetTagsCollectionName(product string) string {
	return product + CollectionSuffixTags
}
