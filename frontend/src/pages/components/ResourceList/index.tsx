import React from "react";
import { type Resource } from "../../../types";

import { ResourceItem } from "./ResourceItem";

import "./ResourceList.scss";

interface ResourceListProps {
  resources: Resource[];
  activeResourceId: string | null;
  onResourceSelect: (resource: Resource) => void;
  isLoading?: boolean;
}

export const ResourceList: React.FC<ResourceListProps> = ({
  resources,
  activeResourceId,
  onResourceSelect,
  isLoading = false,
}) => {
  if (isLoading) {
    return (
      <div className="resource-list-container">
        <div className="resource-list-loading">
          <div className="spinner"></div>
          <p>Loading resources...</p>
        </div>
      </div>
    );
  }

  if (resources.length === 0) {
    return (
      <div className="resource-list-container">
        <div className="resource-list-empty">
          <p>No resources found. Try adjusting your search or filters.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="resource-list-container">
      <div className="resource-list">
        {resources.map((resource) => (
          <ResourceItem
            key={resource.id}
            resource={resource}
            isActive={resource.id === activeResourceId}
            onClick={onResourceSelect}
          />
        ))}
      </div>
    </div>
  );
};
