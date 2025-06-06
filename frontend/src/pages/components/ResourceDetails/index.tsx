import { X, Video, File, ExternalLink, FileText } from "lucide-react";

import { type Resource } from "../../../types";

import "./ResourceDetails.scss";

interface ResourceDetailsProps {
  resource: Partial<Resource> & {
    file?: File | null;
    thumbnail?: File | null;
  };
  onClose: () => void;

  isPreview?: boolean;
}

export const ResourceDetails: React.FC<ResourceDetailsProps> = ({
  resource,
  onClose,
  isPreview = false,
}) => {
  const renderContent = () => {
    switch (resource.type) {
      case "video":
        if (isPreview && resource.file) {
          const videoUrl = URL.createObjectURL(resource.file);
          return (
            <video
              controls
              className="resource-details-video"
              onLoadedData={() => URL.revokeObjectURL(videoUrl)}
            >
              <source src={videoUrl} type={resource.file.type} />
              Your browser does not support the video tag.
            </video>
          );
        }
        return (
          <div className="resource-details-placeholder">
            <Video className="resource-details-placeholder-icon" />
            <p>Video file is not supported</p>
          </div>
        );

      case "pdf":
        if (isPreview && resource.file) {
          const pdfUrl = URL.createObjectURL(resource.file);
          return (
            <iframe
              src={pdfUrl}
              className="resource-details-pdf"
              title="PDF Preview"
              onLoad={() => URL.revokeObjectURL(pdfUrl)}
            />
          );
        }
        return (
          <div className="resource-details-placeholder">
            <File className="resource-details-placeholder-icon" />
            <p>PDF file not supported</p>
          </div>
        );

      case "article":
        if (isPreview && resource.url) {
          return (
            <div className="resource-details-article">
              <div className="resource-details-article-header">
                <ExternalLink />
                <span>External Article</span>
              </div>
              <a
                href={resource.url}
                target="_blank"
                rel="noopener noreferrer"
                className="resource-details-article-link"
              >
                {resource.url.length > 50
                  ? resource.url.substring(0, 50) + "..."
                  : resource.url}
              </a>
              <p className="resource-details-article-note">
                Click the link above to view the article in a new tab
              </p>
            </div>
          );
        }
        return (
          <div className="resource-details-placeholder">
            <ExternalLink className="resource-details-placeholder-icon" />
            <p>Article URL not provided</p>
          </div>
        );

      default:
        return (
          <div className="resource-details-placeholder">
            <FileText className="resource-details-placeholder-icon" />
            <p>{isPreview ? `No preview available` : `No details available`}</p>
          </div>
        );
    }
  };

  return (
    <div className="resource-details">
      <div className="resource-details-header">
        <h3 className="resource-details-title">
          {resource.title || `Resource ${isPreview ? "Preview" : "Detials"}`}
        </h3>
        <button
          onClick={onClose}
          className="resource-details-close"
          aria-label="Close preview"
        >
          <X />
        </button>
      </div>
      <div className="resource-details-content">{renderContent()}</div>
      {resource.description && (
        <div className="resource-details-description">
          <h4>Description:</h4>
          <p>{resource.description}</p>
        </div>
      )}
    </div>
  );
};
