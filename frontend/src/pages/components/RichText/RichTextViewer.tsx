import React from "react";
import "./RichTextViewer.scss";

interface RichTextViewerProps {
  content: string;
  className?: string;
}

export const RichTextViewer: React.FC<RichTextViewerProps> = ({ content, className = "" }) => {
  // Basic sanitization
  const sanitizeHtml = (html: string): string => {
    // Remove script tags and dangerous attributes
    return html
      .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, "")
      .replace(/on\w+="[^"]*"/g, "")
      .replace(/javascript:/gi, "");
  };

  const isEmpty = (html: string): boolean => {
    if (!html) return true;

    // Check if content is empty or just contains empty tags
    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = html;
    const text = tempDiv.textContent || tempDiv.innerText || "";
    return text.trim().length === 0;
  };

  if (isEmpty(content)) {
    return (
      <div className={`richtext-viewer richtext-viewer--empty ${className}`}>
        <span className="richtext-viewer__empty-text">No content available</span>
      </div>
    );
  }

  return (
    <div
      className={`richtext-viewer ${className}`}
      dangerouslySetInnerHTML={{
        __html: sanitizeHtml(content),
      }}
    />
  );
};
