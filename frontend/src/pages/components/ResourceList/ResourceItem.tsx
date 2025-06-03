import React from "react";
import { Video, FileText, Tag } from "lucide-react";

import { type Resource } from "../../../types";

import "./ResourceList.scss";

interface ResourceItemProps {
  resource: Resource;
  isActive: boolean;
  onClick: (resource: Resource) => void;
}

export const ResourceItem: React.FC<ResourceItemProps> = ({ resource, isActive, onClick }) => {
  const { title, description, type, tags } = resource;

  return (
    <div
      className={`resource-item ${isActive ? "active" : ""}`}
      onClick={() => onClick(resource)}
    >
      <div className="resource-icon">{type === "video" ? <Video size={20} /> : <FileText size={20} />}</div>

      <div className="resource-content">
        <h3 className="resource-title">{title}</h3>
        <p className="resource-description">{description}</p>

        <div className="resource-meta">
          {tags && tags.length > 0 && (
            <div className="meta-item">
              <Tag size={14} />
              <span>{tags.slice(0, 2).join(", ")}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
