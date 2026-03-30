package handlers

import (
	contract "learninghub/contract"
	"learninghub/models"
)

func mapResourceToContract(resource models.Resource) contract.Resource {
	var thumbnail *string
	if resource.ThumbnailURL != "" {
		thumbnail = &resource.ThumbnailURL
	}

	return contract.Resource{
		Id:           resource.ID,
		Title:        resource.Title,
		Description:  resource.Description,
		Type:         contract.ResourceType(resource.Type),
		Url:          resource.URL,
		ThumbnailUrl: thumbnail,
		Tags:         resource.Tags,
		CreatedAt:    resource.CreatedAt,
		UpdatedAt:    resource.UpdatedAt,
	}
}

func mapTagToContract(tag models.Tag) contract.Tag {
	return contract.Tag{
		Name:       tag.Name,
		UsageCount: tag.UsageCount,
	}
}
