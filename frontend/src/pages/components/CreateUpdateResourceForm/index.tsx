import React, { useState, useCallback, type JSX, useEffect, useRef } from "react";
import { Upload, Plus, Edit3, Save, Eye, X, Video, File, ExternalLink, FileText, Trash2 } from "lucide-react";

import {
  RESOURCE_TYPES,
  type CreateResourcePayload,
  type Resource,
  type ResourceType,
  type UpdateResourcePayload,
} from "../../../types";

import { SearchSelectInput, type Item } from "../../../components/SearchSelectInput";

import { ResourceDetails } from "../ResourceDetails";

import { useCreateResource, useUpdateResource } from "../../../services/resources";
import { useTags } from "../../../services/tags";

import { usePrevious } from "../../../hooks";

import "./CreateUpdateResourceForm.scss";

interface TFormData extends Partial<Resource> {
  file: File | null;
  thumbnail: File | null;
}

interface CreateUpdateResourceProps {
  onCancel: () => void;
  resource?: Resource;
  onSuccess?: () => void;
}

const defaultType = "video";

export const CreateUpdateResourceForm: React.FC<CreateUpdateResourceProps> = ({ resource, onCancel, onSuccess }) => {
  const [formData, setFormData] = useState<TFormData>({
    title: resource?.title || "",
    description: resource?.description || "",
    type: resource?.type || defaultType,
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
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const prevType = usePrevious(formData.type, defaultType);

  const { data: tags = [], isFetching: isTagsFetching } = useTags();

  const { mutate: createResource, isPending: isCreatingResource } = useCreateResource({
    onSuccess: () => {
      onCancel();
      onSuccess?.();
    },
  });

  const { mutate: updateResource, isPending: isUpdatingResource } = useUpdateResource({
    onSuccess: () => {
      onCancel();
      onSuccess?.();
    },
  });

  const isDisabled = isTagsFetching || isCreatingResource || isUpdatingResource;

  const handleRemoveFile = useCallback((fileType: "file" | "thumbnail") => {
    setFormData((prev) => ({ ...prev, [fileType]: null }));
    // Clear the input using ref
    if (fileType === "file" && fileInputRef.current) {
      fileInputRef.current.value = "";
    }
    if (fileType === "thumbnail" && thumbnailInputRef.current) {
      thumbnailInputRef.current.value = "";
    }
    // remove thumbnail preview
    if (fileType === "thumbnail") {
      setThumbnailPreviewUrl("");
    }
  }, []);

  // reset URL and file state on type chanage when create
  useEffect(() => {
    if (!resource && prevType !== formData.type) {
      setFormData((prev) => ({ ...prev, url: "", file: null }));

      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  }, [resource, prevType, formData.type]);

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

    // update the file and reset url
    setFormData((prev) => ({ ...prev, [fileType]: file, [fileType === "file" ? "url" : "thumbnailUrl"]: "" }));

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

    if (!formData.type?.trim()) {
      errors.type = "Type is required";
    }

    if ((formData.type === "video" || formData.type === "pdf") && !formData.file) {
      if (!resource) {
        errors.file = `Please select a ${formData.type} file`;
      } else if (resource && !formData.url) {
        errors.file = `Please select a ${formData.type} file`;
      }
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

    const payload = getResourcePayload(resource, formData);
    if (resource) {
      updateResource(payload as UpdateResourcePayload);
    } else {
      createResource(payload as CreateResourcePayload);
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
                className={`form-field-input ${validationErrors.title ? "form-field-input-error" : ""}`}
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
                className={`form-field-textarea ${validationErrors.description ? "form-field-textarea-error" : ""}`}
                placeholder="Describe what this resource covers..."
                required
              />
              {validationErrors.description && <span className="form-field-error">{validationErrors.description}</span>}
            </div>

            {/* Type Selection */}
            <div className="form-field">
              <label className="form-field-label">Resource Type *</label>
              <div className="resource-type-selector">
                {Object.values(RESOURCE_TYPES).map((type) => (
                  <button
                    key={type}
                    type="button"
                    onClick={() => {
                      setFormData((prev) => ({ ...prev, type }));
                    }}
                    className={`resource-type-selector-option ${
                      formData.type === type ? "resource-type-selector-option-active" : ""
                    }`}
                    disabled={resource && resource.type !== type}
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
                onSelectedItemsChange={handleTagsChange}
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
                  className={`file-upload ${dragOver ? "file-upload-drag-over" : ""}`}
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
                    ref={fileInputRef}
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
                        disabled={isDisabled}
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
            {(formData.type === "article" || (resource && resource.url && !formData.file)) && (
              <div className="form-field">
                <label className="form-field-label">Resource URL *</label>
                <input
                  type="url"
                  name="url"
                  value={formData.url}
                  onChange={handleInputChange}
                  className={`form-field-input ${validationErrors.url ? "form-field-input-error" : ""}`}
                  placeholder="https://example.com/article"
                  disabled={resource && resource.url && resource.type !== "article" ? true : false}
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
                        handleRemoveFile("thumbnail");
                      }}
                      aria-label="Remove thumbnail"
                      disabled={isDisabled || !!resource}
                    >
                      <X size={16} />
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
          {/* Form Actions */}
          <div className="form-actions">
            <button
              type="button"
              onClick={onCancel}
              className="form-actions-button form-actions-button-secondary"
            >
              Cancel
            </button>

            {canPreview() && (
              <button
                type="button"
                onClick={() => setShowPreview(true)}
                className="form-actions-button form-actions-button-outline"
                disabled={isDisabled}
              >
                <Eye size={16} />
                Preview
              </button>
            )}

            <button
              type="button"
              onClick={handleSubmit}
              className="form-actions-button form-actions-button-primary"
              disabled={isDisabled}
            >
              <Save size={16} />
              {resource ? "Update Resource" : "Create Resource"}
            </button>
          </div>
        </div>
      </div>

      {/* Preview Modal */}
      {showPreview && (
        <div className="create-update-resource-preview-overlay">
          <ResourceDetails
            resource={formData}
            onClose={() => setShowPreview(false)}
            isPreview
          />
        </div>
      )}
    </>
  );
};

function areTagsModified(originalTags: string[] = [], currentTags: string[] = []): boolean {
  if (originalTags.length !== currentTags.length) return true;
  const originalSet = new Set(originalTags);
  const currentSet = new Set(currentTags);
  for (const tag of originalSet) {
    if (!currentSet.has(tag)) return true;
  }
  return false;
}

function getResourcePayload(resource: Resource | undefined, formData: TFormData) {
  if (!resource) {
    // For create, send all fields (except nulls)
    return {
      title: formData.title!,
      description: formData.description!,
      type: formData.type!,
      ...(formData.tags ? { tags: formData.tags.join(",") } : {}),
      ...(formData.url && { url: formData.url }),
      ...(formData.thumbnailUrl && { url: formData.thumbnailUrl }),
      ...(formData.file && { file: formData.file }),
      ...(formData.thumbnail && { thumbnail: formData.thumbnail }),
    };
  }

  // For update, only send changed fields
  const delta: Partial<UpdateResourcePayload> = {};

  if (formData.title !== resource.title) delta.title = formData.title!;
  if (formData.description !== resource.description) delta.description = formData.description!;
  if (formData.type !== resource.type) delta.type = formData.type!;
  if (formData.url !== resource.url) delta.url = formData.url!;
  if (formData.thumbnailUrl !== resource.thumbnailUrl) delta.thumbnailUrl = formData.thumbnailUrl!;
  if (areTagsModified(resource.tags, formData.tags)) delta.tags = formData.tags!.join(",");
  if (formData.file) delta.file = formData.file;
  if (formData.thumbnail) delta.thumbnail = formData.thumbnail;

  delta.id = resource.id;

  return delta;
}
