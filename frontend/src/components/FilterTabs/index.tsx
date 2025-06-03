import React from "react";
import { Video, FileText, Layers } from "lucide-react";
import { type ResourceType } from "../../types";

import "./FilterTabs.scss";

interface FilterTabsProps {
  activeFilter: ResourceType | null;
  onFilterChange: (filter: ResourceType | null) => void;
}

export const FilterTabs: React.FC<FilterTabsProps> = ({ activeFilter, onFilterChange }) => {
  return (
    <div className="filter-tabs">
      <button
        className={`filter-tab ${activeFilter === null ? "active" : ""}`}
        onClick={() => onFilterChange(null)}
        aria-pressed={activeFilter === null}
      >
        <Layers size={18} />
        <span>All</span>
      </button>

      <button
        className={`filter-tab ${activeFilter === "video" ? "active" : ""}`}
        onClick={() => onFilterChange("video")}
        aria-pressed={activeFilter === "video"}
      >
        <Video size={18} />
        <span>Videos</span>
      </button>

      <button
        className={`filter-tab ${activeFilter === "article" ? "active" : ""}`}
        onClick={() => onFilterChange("article")}
        aria-pressed={activeFilter === "article"}
      >
        <FileText size={18} />
        <span>Articles</span>
      </button>
    </div>
  );
};
