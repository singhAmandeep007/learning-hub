import React from "react";
import { FileText, ExternalLink } from "lucide-react";

import { type Resource } from "../../../types";
import "./ContentDisplay.scss";

interface ArticleViewerProps {
  resource: Resource;
}

export const ArticleViewer: React.FC<ArticleViewerProps> = ({ resource }) => {
  return (
    <div className="article-viewer">
      <div className="article-header">
        <h2>{resource.title}</h2>
        <div className="article-meta">
          {resource.tags.length > 0 && (
            <div className="article-tags">
              {resource.tags.map((tag, index) => (
                <span
                  key={index}
                  className="tag"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="article-container">
        {/* This is a placeholder for the PDF embed */}
        <div className="article-placeholder">
          <div className="placeholder-content">
            <FileText size={48} />
            <p>PDF Document: {resource.title}</p>
            <p className="placeholder-note">
              In the actual implementation, this would be a PDF viewer or iframe embed.
            </p>

            <a
              href={resource.url}
              target="_blank"
              rel="noopener noreferrer"
              className="view-original-button"
            >
              <span>View Original Article</span>
              <ExternalLink size={16} />
            </a>
          </div>
        </div>
      </div>

      <div className="article-description">
        <h3>Description</h3>
        <p>{resource.description}</p>
      </div>
    </div>
  );
};
