import { X, Video, File, ExternalLink, FileText } from "lucide-react";

import { type Resource } from "../../../types";

import { RichTextViewer } from "../RichText";

import "./ResourceDetails.scss";

interface ResourceDetailsProps {
  resource: Partial<Resource> & {
    file?: File | null;
    thumbnail?: File | null;
  };
  onClose: () => void;

  isPreview?: boolean;
}

export const ResourceDetails: React.FC<ResourceDetailsProps> = ({ resource, onClose, isPreview = false }) => {
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
              <source
                src={videoUrl}
                type={resource.file.type}
              />
              Your browser does not support the video tag.
            </video>
          );
        }
        if (!isPreview && resource.url) {
          return (
            <video
              controls
              className="resource-details-video"
            >
              <source
                src={resource.url}
                type="video/mp4"
              />
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
            <object
              data={pdfUrl}
              type="application/pdf"
              className="resource-details-pdf"
              onLoad={() => URL.revokeObjectURL(pdfUrl)}
            >
              <div className="resource-details-pdf-fallback">
                <a
                  href={pdfUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="resource-details-pdf-fallback-link"
                >
                  <ExternalLink size={16} />
                  <span>External PDF</span>
                </a>
                <p className="resource-details-pdf-fallback-note">Click the link above to view the PDF in a new tab</p>
              </div>
            </object>
          );
        }
        if (!isPreview && resource.url) {
          return (
            <object
              data={resource.url}
              type="application/pdf"
              className="resource-details-pdf"
            >
              <div className="resource-details-pdf-fallback">
                <a
                  href={resource.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="resource-details-pdf-fallback-link"
                >
                  <ExternalLink size={16} />
                  <span>External PDF</span>
                </a>
                <p className="resource-details-pdf-fallback-note">Click the link above to view the PDF in a new tab</p>
              </div>
            </object>
          );
        }
        return (
          <div className="resource-details-placeholder">
            <File className="resource-details-placeholder-icon" />
            <p>PDF file not supported</p>
          </div>
        );

      case "article":
        if (resource.url) {
          return (
            <div className="resource-details-article">
              <a
                href={resource.url}
                target="_blank"
                rel="noopener noreferrer"
                className="resource-details-article-link"
              >
                <ExternalLink size={16} />
                <span>External Article</span>
              </a>
              <p className="resource-details-article-note">Click the link above to view the article in a new tab</p>
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
        <h3 className="resource-details-title">{resource.title || `Resource ${isPreview ? "Preview" : "Details"}`}</h3>
        <button
          onClick={onClose}
          className="resource-details-close"
          aria-label="Close preview"
        >
          <X size={16} />
        </button>
      </div>
      <div className="resource-details-content">{renderContent()}</div>
      {resource.description && (
        <div className="resource-details-description">
          <h4>Description:</h4>

          <RichTextViewer content={resource.description} />
        </div>
      )}
    </div>
  );
};
