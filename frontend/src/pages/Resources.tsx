import { useCallback, useEffect, useState } from "react";

import { Plus, Search, X, FileText } from "lucide-react";

import { ResourceCard } from "./components/ResourceCard";
import { ScrollableTags } from "./components/ScrollableTags";
import { CreateUpdateResourceForm } from "./components/CreateUpdateResourceForm";

import { useTags } from "../services/tags";
import { useDeleteResource, useResources } from "../services/resources";

import { type Resource, type ResourcesFilters, type Tag } from "../types";

import "./Resources.scss";

export const Resources = () => {
  const [search, setSearch] = useState("");
  const [selectedTags, setSelectedTags] = useState<Tag["name"][]>([]);
  const [selectedType, setSelectedType] = useState<ResourcesFilters["type"]>("all");

  const [queryParams, setQueryParams] = useState(() => ({
    search: "",
    type: "all" as ResourcesFilters["type"],
    tags: [] as Tag["name"][],
  }));

  const [editingResource, setEditingResource] = useState<Resource | null>(null);

  // Initialize from URL parameters
  useEffect(() => {
    const urlParams = getUrlParams();

    setSearch(urlParams.search);
    setSelectedTags(urlParams.tags);
    setSelectedType(urlParams.type);

    setQueryParams({
      search: urlParams.search,
      type: urlParams.type,
      tags: urlParams.tags,
    });
  }, []);

  const [showCreateForm, setShowCreateForm] = useState(false);

  const {
    data: resources = {
      data: [],
      hasMore: false,
    },
  } = useResources({
    ...(queryParams.search ? { search: queryParams.search } : {}),
    ...(queryParams.type && queryParams.type !== "all" ? { type: queryParams.type } : {}),
    ...(queryParams.tags.length > 0 ? { tags: queryParams.tags } : {}),
  });

  const { data: tags = [], isFetching: isFetchingTags } = useTags();

  const { mutate: deleteResource, isPending: isDeletingResource } = useDeleteResource();

  const isSearchDisabled = isFetchingTags || isDeletingResource;

  const handleSearch = useCallback(() => {
    updateUrlParams(search, selectedTags, selectedType);
    setQueryParams({
      search,
      type: selectedType,
      tags: selectedTags,
    });
  }, [search, selectedTags, selectedType]);

  return (
    <div className="resources">
      {/* Header */}
      <div className="resources-header">
        <div className="resources-header-container">
          <div className="resources-header-content">
            <div className="resources-header-info">
              <h1 className="resources-title">Learning Hub</h1>
              <p className="resources-subtitle">
                Find tutorials, guides, and resources to help you get the most out of the Migration App.
              </p>
            </div>
            <button
              onClick={() => setShowCreateForm(true)}
              className="resources-create-btn"
            >
              <Plus
                className="resources-create-btn-icon"
                size={16}
              />
              Create
            </button>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="resources-main">
        <div className="resources-filters">
          <div className="resources-filters-content">
            {/* Search */}
            <div className="resources-search">
              <div className="resources-search-wrapper">
                <Search
                  className="resources-search-icon"
                  size={16}
                />
                <input
                  type="text"
                  placeholder="Search resources..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="resources-search-input"
                />
                {/* Type Filter */}
                <div className="resources-type-filter">
                  <select
                    value={selectedType}
                    onChange={(e) => setSelectedType(e.target.value as ResourcesFilters["type"])}
                    className="resources-type-select"
                  >
                    <option value="all">All Types</option>
                    <option value="video">Videos</option>
                    <option value="pdf">PDFs</option>
                    <option value="article">Articles</option>
                  </select>
                </div>

                <button
                  className="resources-search-btn"
                  onClick={handleSearch}
                  disabled={isSearchDisabled}
                >
                  Search
                </button>
              </div>
            </div>

            <div className="resources-tag-filters-container">
              {/* Tag Filter */}
              <ScrollableTags
                tags={tags}
                selectedTags={selectedTags}
                setSelectedTags={setSelectedTags}
              />

              {/* Clear Filters */}
              {(search || selectedType !== "all" || selectedTags.length > 0) && (
                <button
                  onClick={() => {
                    setSearch("");
                    setSelectedType("all");
                    setSelectedTags([]);
                  }}
                  className="resources-clear-filters"
                >
                  <X
                    className="resources-clear-filters-icon"
                    size={12}
                  />
                  Clear
                </button>
              )}
            </div>
          </div>
        </div>

        {/* Results Count */}
        <div className="resources-results-count">
          <p className="resources-results-text">
            Showing {resources.data.length} of {resources.data.length} resources
          </p>
        </div>

        {/* Resources Grid */}
        <div className="resources-grid">
          {resources.data.map((resource) => (
            <ResourceCard
              key={resource.id}
              resource={resource}
              onEdit={setEditingResource}
              onDelete={(id) => {
                deleteResource({ id });
              }}
            />
          ))}
        </div>

        {resources.data.length === 0 && (
          <div className="resources-empty-state">
            <FileText className="resources-empty-state-icon" />
            <h3 className="resources-empty-state-title">No resources found</h3>
            <p className="resources-empty-state-text">
              {search || selectedType !== "all" || selectedTags.length > 0
                ? "Try adjusting your filters or search query"
                : "Get started by creating your first resource"}
            </p>
          </div>
        )}
      </div>

      {showCreateForm && <CreateUpdateResourceForm onCancel={() => setShowCreateForm(false)} />}

      {editingResource && (
        <CreateUpdateResourceForm
          onCancel={() => setEditingResource(null)}
          resource={editingResource}
        />
      )}
    </div>
  );
};

// URL parameter utilities
const getUrlParams = () => {
  const params = new URLSearchParams(window.location.search);
  return {
    search: params.get("search") || "",
    tags: params.get("tags") && typeof params.get("tags") === "string" ? params.get("tags")!.split(",") : [],
    type: params.get("type") || "all",
  } as Required<ResourcesFilters>;
};

const updateUrlParams = (
  search: ResourcesFilters["search"],
  tags: ResourcesFilters["tags"],
  type: ResourcesFilters["type"]
) => {
  const params = new URLSearchParams();

  if (search) params.set("search", search);
  if (type && type !== "all") params.set("type", type);
  if (Array.isArray(tags) && tags.length > 0) params.set("tags", tags.join(","));

  const newUrl = `${window.location.pathname}${params.toString() ? "?" + params.toString() : ""}`;
  window.history.pushState({}, "", newUrl);
};
