import { useState, useCallback } from "react";
import { Video, File, ExternalLink, FileText, Eye, Edit3, Trash2, Tag } from "lucide-react";

import { type Resource, RESOURCE_TYPES, type ResourceType } from "../../../types";

import "./ResourceCard.scss";
import { ResourceDetails } from "../ResourceDetails";

interface ResourceCardProps {
  resource: Resource;
  onEdit: (resource: Resource) => void;
  onDelete: (id: Resource["id"]) => void;
}

export const ResourceCard = ({ resource, onEdit, onDelete }: ResourceCardProps) => {
  const [showDetails, setShowDetails] = useState<boolean>(false);

  const getTypeIcon = useCallback((type: ResourceType) => {
    switch (type) {
      case RESOURCE_TYPES.video:
        return <Video className="resource-card-type-icon resource-card-type-icon-video" />;
      case RESOURCE_TYPES.pdf:
        return <File className="resource-card-type-icon resource-card-type-icon-pdf" />;
      case RESOURCE_TYPES.article:
        return <ExternalLink className="resource-card-type-icon resource-card-type-icon-article" />;
      default:
        return <FileText className="resource-card-type-icon" />;
    }
  }, []);

  const handleShowDetails = useCallback(() => {
    setShowDetails(true);
  }, []);

  const handleCloseDetails = useCallback(() => {
    setShowDetails(false);
  }, []);

  const handleEdit = useCallback(() => {
    onEdit(resource);
  }, [onEdit, resource]);

  const handleDelete = useCallback(() => {
    if (window.confirm("Are you sure you want to delete this resource?")) {
      onDelete(resource.id);
    }
  }, [onDelete, resource.id]);

  return (
    <>
      <div className="resource-card">
        <div className="resource-card-header">
          <div className="resource-card-type">
            {getTypeIcon(resource.type)}
            <span className="resource-card-type-label">{resource.type}</span>
          </div>
          <div className="resource-card-actions">
            <button
              onClick={handleShowDetails}
              className="resource-card-action-btn resource-card-action-btn-view"
              title="View Resource"
              type="button"
            >
              <Eye className="resource-card-action-icon" />
            </button>
            <button
              onClick={handleEdit}
              className="resource-card-action-btn resource-card-action-btn-edit"
              title="Edit Resource"
              type="button"
            >
              <Edit3 className="resource-card-action-icon" />
            </button>
            <button
              onClick={handleDelete}
              className="resource-card-action-btn resource-card-action-btn-delete"
              title="Delete Resource"
              type="button"
            >
              <Trash2 className="resource-card-action-icon" />
            </button>
          </div>
        </div>

        {resource.thumbnailUrl && (
          <div className="resource-card-thumbnail">
            <img
              src={resource.thumbnailUrl}
              alt="Thumbnail"
              className="resource-card-thumbnail-image"
            />
          </div>
        )}

        {!resource.thumbnailUrl && (
          <div className="resource-card-thumbnail  resource-card-thumbnail-fallback">{getTypeIcon(resource.type)}</div>
        )}

        <h3 className="resource-card-title">{resource.title}</h3>
        <p className="resource-card-description">{resource.description}</p>

        <div className="resource-card-tags">
          {resource.tags.map((tag) => (
            <span
              key={tag}
              className="resource-card-tag"
            >
              <Tag className="resource-card-tag-icon" />
              {tag}
            </span>
          ))}
        </div>

        <div className="resource-card-dates">
          Created: {new Date(resource.createdAt).toLocaleDateString()}
          {resource.updatedAt !== resource.createdAt && (
            <span> • Updated: {new Date(resource.updatedAt).toLocaleDateString()}</span>
          )}
        </div>
      </div>

      {showDetails && (
        <div className="create-update-resource-preview-overlay">
          <ResourceDetails
            resource={resource}
            onClose={handleCloseDetails}
          />
        </div>
      )}
    </>
  );
};
