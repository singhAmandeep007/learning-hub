import React, { useState, useCallback, type JSX, useEffect, useRef } from "react";
import { Upload, Plus, Edit3, Save, Eye, X, Video, File, ExternalLink, FileText, Trash2 } from "lucide-react";

import "./CreateUpdateResourceForm.scss";

import {
  RESOURCE_TYPES,
  type CreateResourcePayload,
  type Resource,
  type ResourceType,
  type UpdateResourcePayload,
} from "../../../types";

import { SearchSelectInput, type Item } from "../../../components/SearchSelectInput";

import { useCreateResource } from "../../../services/resources";

import { useTags } from "../../../services/tags";

interface FormData extends Partial<Resource> {
  file: File | null;
  thumbnail: File | null;
}

const ResourcePreview: React.FC<{
  resource: FormData;
  onClose: () => void;
}> = ({ resource, onClose }) => {
  const renderPreviewContent = () => {
    switch (resource.type) {
      case "video":
        if (resource.file) {
          const videoUrl = URL.createObjectURL(resource.file);
          return (
            <video
              controls
              className="resource-preview-video"
              onLoadedData={() => URL.revokeObjectURL(videoUrl)}
            >
              <source
                src={videoUrl}
                type={resource.file.type}
              />
              Your browser does not support the video tag.
            </video>
          );
        }
        return (
          <div className="resource-preview-placeholder">
            <Video className="resource-preview-placeholder-icon" />
            <p>Video file not available for preview</p>
          </div>
        );

      case "pdf":
        if (resource.file) {
          const pdfUrl = URL.createObjectURL(resource.file);
          return (
            <iframe
              src={pdfUrl}
              className="resource-preview-pdf"
              title="PDF Preview"
              onLoad={() => URL.revokeObjectURL(pdfUrl)}
            />
          );
        }
        return (
          <div className="resource-preview-placeholder">
            <File className="resource-preview-placeholder-icon" />
            <p>PDF file not available for preview</p>
          </div>
        );

      case "article":
        if (resource.url) {
          return (
            <div className="resource-preview-article">
              <div className="resource-preview-article-header">
                <ExternalLink />
                <span>External Article</span>
              </div>
              <a
                href={resource.url}
                target="_blank"
                rel="noopener noreferrer"
                className="resource-preview-article-link"
              >
                {resource.url.length > 50 ? resource.url.substring(0, 50) + "..." : resource.url}
              </a>
              <p className="resource-preview-article-note">Click the link above to view the article in a new tab</p>
            </div>
          );
        }
        return (
          <div className="resource-preview-placeholder">
            <ExternalLink className="resource-preview-placeholder-icon" />
            <p>Article URL not provided</p>
          </div>
        );

      default:
        return (
          <div className="resource-preview-placeholder">
            <FileText className="resource-preview-placeholder-icon" />
            <p>No preview available</p>
          </div>
        );
    }
  };

  return (
    <div className="resource-preview">
      <div className="resource-preview-header">
        <h3 className="resource-preview-title">{resource.title || "Resource Preview"}</h3>
        <button
          onClick={onClose}
          className="resource-preview-close"
          aria-label="Close preview"
        >
          <X />
        </button>
      </div>
      <div className="resource-preview-content">{renderPreviewContent()}</div>
      {resource.description && (
        <div className="resource-preview-description">
          <h4>Description:</h4>
          <p>{resource.description}</p>
        </div>
      )}
    </div>
  );
};

interface CreateUpdateResourceProps {
  resource?: Resource;
  onSuccess?: () => void;
  onCancel: () => void;
}

export const CreateUpdateResourceForm: React.FC<CreateUpdateResourceProps> = ({ resource, onCancel }) => {
  const [formData, setFormData] = useState<FormData>({
    title: resource?.title || "",
    description: resource?.description || "",
    type: resource?.type || "video",
    url: resource?.url || "",
    thumbnailUrl: resource?.thumbnailUrl || "",
    tags: resource?.tags || [],
    file: null,
    thumbnail: null,
  });

  const [dragOver, setDragOver] = useState<boolean>(false);
  const [thumbnailPreviewUrl, setThumbnailPreviewUrl] = useState<string>(resource?.thumbnailUrl || "");
  const [showPreview, setShowPreview] = useState<boolean>(false);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  const thumbnailInputRef = useRef<HTMLInputElement | null>(null);

  const { data: tags } = useTags();

  const { mutate: createResource } = useCreateResource();

  // Effect to handle URL input state based on file attachment
  useEffect(() => {
    if (formData.file && (formData.type === "video" || formData.type === "pdf")) {
      setFormData((prev) => ({ ...prev, url: "" }));
    }
  }, [formData.file, formData.type]);

  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
      const { name, value } = e.target;
      setFormData((prev) => ({ ...prev, [name]: value }));

      // Clear validation error when user starts typing
      if (validationErrors[name]) {
        setValidationErrors((prev) => ({ ...prev, [name]: "" }));
      }
    },
    [validationErrors]
  );

  const handleFileChange = useCallback((e: React.ChangeEvent<HTMLInputElement>, fileType: "file" | "thumbnail") => {
    const file = e.target.files?.[0];
    if (!file) return;

    setFormData((prev) => ({ ...prev, [fileType]: file }));

    if (fileType === "thumbnail") {
      const reader = new FileReader();
      reader.onload = (event) => {
        if (event.target?.result) {
          setThumbnailPreviewUrl(event.target.result as string);
        }
      };
      reader.readAsDataURL(file);
    }
  }, []);

  const handleRemoveFile = useCallback((fileType: "file" | "thumbnail") => {
    setFormData((prev) => ({ ...prev, [fileType]: null }));
  }, []);

  const handleDragOver = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setDragOver(true);
  }, []);

  const handleDragLeave = useCallback(() => {
    setDragOver(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setDragOver(false);
    const file = e.dataTransfer.files[0];
    if (file) {
      setFormData((prev) => ({ ...prev, file }));
    }
  }, []);

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.title?.trim()) {
      errors.title = "Title is required";
    }

    if (!formData.description?.trim()) {
      errors.description = "Description is required";
    }

    if ((formData.type === "video" || formData.type === "pdf") && !formData.file && !resource) {
      errors.file = `Please select a ${formData.type} file`;
    }

    if (formData.type === "article" && !formData.url?.trim()) {
      errors.url = "URL is required for articles";
    }

    if (formData.tags?.length === 0) {
      errors.tags = "At least one tag is required";
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = useCallback(() => {
    if (!validateForm()) {
      return;
    }

    const payload: CreateResourcePayload = {
      title: formData.title!,
      description: formData.description!,
      type: formData.type!,
      tags: formData.tags!.join(",") || "",
      ...(formData.url && { url: formData.url }),
      ...(formData.thumbnailUrl && { url: formData.thumbnailUrl }),
      ...(formData.file && { file: formData.file }),
      ...(formData.thumbnail && { thumbnail: formData.thumbnail }),
    };

    if (resource) {
      (payload as UpdateResourcePayload).id = resource.id; // Include ID for update
      console.log("Updating resource with payload:", payload);
    } else {
      createResource(payload);
    }
    // eslint-disable-next-line
  }, [formData, resource, createResource]);

  const handleTagsChange = useCallback((selectedItems: Item[]) => {
    setFormData((prev) => ({
      ...prev,
      tags: selectedItems.map((item) => item.name),
    }));
  }, []);

  const getTypeIcon = (type: ResourceType): JSX.Element => {
    switch (type) {
      case "video":
        return <Video />;
      case "pdf":
        return <File />;
      case "article":
        return <ExternalLink />;
      default:
        return <FileText />;
    }
  };

  const canPreview = (): boolean => {
    return !!(
      formData.title &&
      formData.description &&
      ((formData.type === "article" && formData.url) ||
        (["pdf", "video"].includes(formData?.type || "") && formData.file))
    );
  };

  return (
    <>
      <div className="create-update-resource-form">
        <div className="create-update-resource-form-container">
          <div className="create-update-resource-form-header">
            <h2 className="create-update-resource-form-title">
              {resource ? <Edit3 /> : <Plus />}
              {resource ? "Edit Resource" : "Create New Resource"}
            </h2>
          </div>

          <div className="create-update-resource-form-content">
            {/* Title */}
            <div className="form-field">
              <label className="form-field-label">Title *</label>
              <input
                type="text"
                name="title"
                value={formData.title}
                onChange={handleInputChange}
                className={`form-field-input ${validationErrors.title ? "form-field-input--error" : ""}`}
                placeholder="Enter resource title..."
                required
              />
              {validationErrors.title && <span className="form-field-error">{validationErrors.title}</span>}
            </div>

            {/* Description */}
            <div className="form-field">
              <label className="form-field-label">Description *</label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                rows={4}
                className={`form-field-textarea ${validationErrors.description ? "form-field-textarea--error" : ""}`}
                placeholder="Describe what this resource covers..."
                required
              />
              {validationErrors.description && <span className="form-field-error">{validationErrors.description}</span>}
            </div>

            {/* Type Selection */}
            <div className="form-field">
              <label className="form-field-label">Resource Type *</label>
              <div className="resource-type-selector">
                {RESOURCE_TYPES.map((type) => (
                  <button
                    key={type}
                    type="button"
                    onClick={() => {
                      setFormData((prev) => ({ ...prev, type }));
                      handleRemoveFile("file");
                    }}
                    className={`resource-type-selector-option ${
                      formData.type === type ? "resource-type-selector-option--active" : ""
                    }`}
                  >
                    {getTypeIcon(type)}
                    <span className="resource-type-selector-label">{type}</span>
                  </button>
                ))}
              </div>
              {validationErrors.type && <span className="form-field-error">{validationErrors.type}</span>}
            </div>

            {/* Tags */}
            <div className="form-field">
              <label className="form-field-label">Tags *</label>
              <SearchSelectInput
                items={tags?.map((tag) => ({ id: tag.name, name: tag.name })) || []}
                placeholder="Search and select tags..."
                onSelectedItemsChange={(tags) => {
                  handleTagsChange(tags);
                  console.log("Selected tags:", tags);
                }}
                initialSelectedItems={formData.tags?.map((tag) => ({ id: tag, name: tag })) || []}
                allowNewTags
              />
              {validationErrors.tags && <span className="form-field-error">{validationErrors.tags}</span>}
            </div>

            {/* File Upload (for video/pdf) */}
            {(formData.type === "video" || formData.type === "pdf") && (
              <div className="form-field">
                <label className="form-field-label">
                  {formData.type === "video" ? "Video File" : "PDF File"} {!resource && "*"}
                </label>
                <div
                  className={`file-upload ${dragOver ? "file-upload--drag-over" : ""}`}
                  onDragOver={handleDragOver}
                  onDragLeave={handleDragLeave}
                  onDrop={handleDrop}
                >
                  <Upload className="file-upload-icon" />
                  <p className="file-upload-text">Drag & drop your {formData.type} file here, or click to browse</p>
                  <input
                    type="file"
                    accept={formData.type === "video" ? "video/*" : ".pdf"}
                    onChange={(e) => handleFileChange(e, "file")}
                    className="file-upload-input"
                    id="file-upload"
                    required={!resource}
                  />
                  <label
                    htmlFor="file-upload"
                    className="file-upload-button"
                  >
                    Choose File
                  </label>

                  {formData.file && (
                    <div className="file-upload-selected-info">
                      <span>Selected: {formData.file.name}</span>
                      <button
                        onClick={() => handleRemoveFile("file")}
                        className="file-upload-remove-button"
                        aria-label="Remove file"
                      >
                        <Trash2 />
                      </button>
                    </div>
                  )}
                </div>
                {validationErrors.file && <span className="form-field-error">{validationErrors.file}</span>}
              </div>
            )}

            {/* URL (for articles or as fallback) */}
            {formData.type === "article" && (
              <div className="form-field">
                <label className="form-field-label">External URL *</label>
                <input
                  type="url"
                  name="url"
                  value={formData.url}
                  onChange={handleInputChange}
                  className={`form-field-input ${validationErrors.url ? "form-field-input--error" : ""}`}
                  placeholder="https://example.com/article"
                />
                {validationErrors.url && <span className="form-field-error">{validationErrors.url}</span>}
              </div>
            )}

            {/* Thumbnail Upload */}
            <div className="form-field">
              <label className="form-field-label">Thumbnail (Optional)</label>
              <div className="thumbnail-upload">
                <div className="thumbnail-upload-input-container">
                  <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => handleFileChange(e, "thumbnail")}
                    className="thumbnail-upload-input"
                    ref={thumbnailInputRef}
                  />
                </div>
                {thumbnailPreviewUrl && (
                  <div className="thumbnail-upload-preview">
                    <img
                      src={thumbnailPreviewUrl}
                      alt="Thumbnail preview"
                      className="thumbnail-upload-image"
                    />
                    <button
                      className="thumbnail-upload-remove"
                      onClick={() => {
                        setThumbnailPreviewUrl("");
                        setFormData((prev) => ({ ...prev, thumbnail: null }));
                        if (thumbnailInputRef.current) {
                          thumbnailInputRef.current.value = ""; // Reset the file input
                        }
                      }}
                      aria-label="Remove thumbnail"
                    >
                      <X size={16} />
                    </button>
                  </div>
                )}
              </div>
            </div>

            {/* Form Actions */}
            <div className="form-actions">
              <button
                type="button"
                onClick={onCancel}
                className="form-actions-button form-actions-button--secondary"
              >
                Cancel
              </button>

              {canPreview() && (
                <button
                  type="button"
                  onClick={() => setShowPreview(true)}
                  className="form-actions-button form-actions-button--outline"
                >
                  <Eye />
                  Preview
                </button>
              )}

              <button
                type="button"
                onClick={handleSubmit}
                className="form-actions-button form-actions-button--primary"
              >
                <Save />
                {resource ? "Update Resource" : "Create Resource"}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Preview Modal */}
      {showPreview && (
        <div className="create-update-resource-preview-overlay">
          <ResourcePreview
            resource={formData}
            onClose={() => setShowPreview(false)}
          />
        </div>
      )}
    </>
  );
};
