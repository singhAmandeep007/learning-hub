import { useCallback, useEffect, useState } from "react";

import { Plus, Search, X, FileText, ChevronRight, ChevronLeft } from "lucide-react";

import { ResourceCard } from "./components/ResourceCard";
import { ScrollableTags } from "./components/ScrollableTags";
import { CreateUpdateResourceForm } from "./components/CreateUpdateResourceForm";

import { useTags } from "../services/tags";
import { useDeleteResource, useResources } from "../services/resources";

import { type Resource, type ResourcesFilters, type Tag } from "../types";

import "./Resources.scss";

export const Resources = () => {
  // filter state
  const [search, setSearch] = useState("");
  const [selectedTags, setSelectedTags] = useState<Tag["name"][]>([]);
  const [selectedType, setSelectedType] = useState<ResourcesFilters["type"]>("all");

  const [editingResource, setEditingResource] = useState<Resource | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  // Pagination state
  const [currentPage, setCurrentPage] = useState(0);
  const limit = 20;

  const calculateCursor = (page: number, itemsPerPage: number) => {
    if (page <= 1) return "";
    return String((page - 1) * itemsPerPage);
  };

  const queryParams = {
    ...(search ? { search } : {}),
    ...(selectedType && selectedType !== "all" ? { type: selectedType } : {}),
    ...(selectedTags.length > 0 ? { tags: selectedTags } : {}),
    cursor: calculateCursor(currentPage, limit),
    limit: String(limit),
  };

  const {
    data: resources = {
      data: [],
      hasMore: false,
      nextCursor: "",
    },
    isFetching: isFetchingResources,
    isLoading: isLoadingResources,
    refetch,
  } = useResources(queryParams);

  const { data: tags = [], isFetching: isFetchingTags } = useTags();

  const { mutate: deleteResource, isPending: isDeletingResource } = useDeleteResource();

  const isSearchDisabled = isFetchingTags || isDeletingResource || isFetchingResources;

  const handleSearch = useCallback(() => {
    // Reset pagination when searching
    setCurrentPage(1);
    refetch();
  }, [refetch]);

  const handleClearFilters = useCallback(() => {
    setSearch("");
    setSelectedType("all");
    setSelectedTags([]);

    setCurrentPage(1);
  }, []);

  const handleNextPage = useCallback(() => {
    if (resources.hasMore) {
      setCurrentPage((prev) => prev + 1);
    }
  }, [resources.hasMore]);

  const handlePrevPage = useCallback(() => {
    if (currentPage > 1) {
      setCurrentPage((prev) => prev - 1);
    }
  }, [currentPage]);

  // Reset pagination when filters change
  useEffect(() => {
    setCurrentPage(1);
  }, [search, selectedTags, selectedType]);

  useEffect(() => {
    refetch();
  }, [currentPage, refetch]);

  const handleResourceChange = useCallback(() => {
    handleClearFilters();
    refetch();
    // refetch will happen automatically due to state changes
  }, [handleClearFilters, refetch]);

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
                  onKeyDown={(e) => e.key === "Enter" && handleSearch()}
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
                  onClick={handleClearFilters}
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
        <div className="resources-results-header">
          <div className="resources-results-count">
            <p className="resources-results-text">
              {isLoadingResources
                ? "Loading resources..."
                : `Showing ${resources.data.length} resources (Page ${currentPage})`}
            </p>
          </div>
          {/* Pagination Controls */}
          {!isLoadingResources && (currentPage > 1 || resources.hasMore) && (
            <div className="resources-pagination">
              <button
                disabled={currentPage <= 1 || isFetchingResources}
                className="resources-pagination-btn resources-pagination-prev"
                onClick={handlePrevPage}
              >
                <ChevronLeft size={16} />
              </button>

              <span className="resources-pagination-info">Page {currentPage}</span>

              <button
                disabled={!resources.hasMore || isFetchingResources}
                className="resources-pagination-btn resources-pagination-next"
                onClick={handleNextPage}
              >
                <ChevronRight size={16} />
              </button>
            </div>
          )}
        </div>

        {/* Empty State */}
        {!isLoadingResources && resources.data.length === 0 && (
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

        {/* Resources Grid */}
        {resources.data.length > 0 && (
          <div className="resources-grid">
            <div className="resources-grid-container">
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
          </div>
        )}
      </div>

      {showCreateForm && (
        <CreateUpdateResourceForm
          onCancel={() => setShowCreateForm(false)}
          onSuccess={handleResourceChange}
        />
      )}

      {editingResource && (
        <CreateUpdateResourceForm
          onCancel={() => setEditingResource(null)}
          resource={editingResource}
          onSuccess={handleResourceChange}
        />
      )}
    </div>
  );
};

// URL parameter utilities
// const getUrlParams = () => {
//   const params = new URLSearchParams(window.location.search);
//   return {
//     search: params.get("search") || "",
//     tags: params.get("tags") && typeof params.get("tags") === "string" ? params.get("tags")!.split(",") : [],
//     type: params.get("type") || "all",
//   } as Required<ResourcesFilters>;
// };

// const updateUrlParams = (
//   search: ResourcesFilters["search"],
//   tags: ResourcesFilters["tags"],
//   type: ResourcesFilters["type"]
// ) => {
//   const params = new URLSearchParams();

//   if (search) params.set("search", search);
//   if (type && type !== "all") params.set("type", type);
//   if (Array.isArray(tags) && tags.length > 0) params.set("tags", tags.join(","));

//   const newUrl = `${window.location.pathname}${params.toString() ? "?" + params.toString() : ""}`;
//   window.history.pushState({}, "", newUrl);
// };
