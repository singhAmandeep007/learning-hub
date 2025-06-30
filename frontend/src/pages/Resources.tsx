import { useCallback, useMemo, useState } from "react";
import { Plus, Search, X, FileText, ChevronRight, ChevronLeft } from "lucide-react";

import { ResourceCard } from "./components/ResourceCard";
import { ScrollableTags } from "./components/ScrollableTags";
import { CreateUpdateResourceForm } from "./components/CreateUpdateResourceForm";

import { useTags } from "../services/tags";
import { useDeleteResource, useResources } from "../services/resources";
import { useResourceFilters } from "../hooks";

import { type Resource } from "../types";

import "./Resources.scss";

export const Resources = () => {
  const { data: tags = [], isFetching: isFetchingTags, isSuccess: hasFetchedTags } = useTags();

  const loadedTags = useMemo(() => tags.map((tag) => tag.name), [tags]);

  const {
    searchInput,
    setSearchInput,

    selectedTags,
    selectedType,
    currentPage,

    setSelectedTags,
    setSelectedType,
    setCurrentPage,
    setActiveSearch,
    handleClearFilters,

    queryParams,
    hasActiveFilters,
  } = useResourceFilters({
    loadedTags,
    hasFetchedTags,
  });

  const [editingResource, setEditingResource] = useState<Resource | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  const {
    data: resources = {
      data: [],
      hasMore: false,
      nextCursor: "",
    },
    isFetching: isFetchingResources,
    isLoading: isLoadingResources,
  } = useResources(queryParams);

  const { mutate: deleteResource, isPending: isDeletingResource } = useDeleteResource();

  const isSearchDisabled = isFetchingTags || isDeletingResource || isFetchingResources;

  const handleNextPage = useCallback(() => {
    if (resources.hasMore) {
      setCurrentPage(currentPage + 1);
    }
  }, [resources.hasMore, currentPage, setCurrentPage]);

  const handlePrevPage = useCallback(() => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  }, [currentPage, setCurrentPage]);

  const handleEditResource = useCallback((resource: Resource) => {
    setEditingResource(resource);
  }, []);

  const handleDeleteResource = useCallback(
    (id: Resource["id"]) => {
      deleteResource({ id });
    },
    [deleteResource]
  );

  const handleSearch = useCallback(() => {
    setActiveSearch(searchInput);
    // Reset pagination when searching
    setCurrentPage(1);
  }, [setActiveSearch, searchInput, setCurrentPage]);

  // Handle Enter key in search input
  const handleSearchKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === "Enter") {
        handleSearch();
      }
    },
    [handleSearch]
  );

  return (
    <div className="resources">
      {/* Header */}
      <div className="resources-header">
        <div className="resources-header-container">
          <div className="resources-header-content">
            <div className="resources-header-info">
              <h1 className="resources-title">Learning Hub</h1>
              <p className="resources-subtitle">
                Find tutorials, guides, and resources to help you get the most out of the App.
              </p>
            </div>
            <button
              onClick={() => setShowCreateForm(true)}
              className="resources-create-btn"
              type="button"
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

      {/* Main Content */}
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
                  value={searchInput}
                  onChange={(e) => setSearchInput(e.target.value)}
                  onKeyDown={handleSearchKeyDown}
                  className="resources-search-input"
                  id="resources-search-input"
                  name="search-query"
                />
                {/* Type Filter */}
                <div className="resources-type-filter">
                  <select
                    value={selectedType}
                    onChange={(e) => setSelectedType(e.target.value as typeof selectedType)}
                    className="resources-type-select"
                    id="resources-type"
                    name="resources-type"
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
                  type="button"
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
              {hasActiveFilters && (
                <button
                  onClick={handleClearFilters}
                  className="resources-clear-filters"
                  type="button"
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
                type="button"
              >
                <ChevronLeft size={16} />
              </button>

              <span className="resources-pagination-info">Page {currentPage}</span>

              <button
                disabled={!resources.hasMore || isFetchingResources}
                className="resources-pagination-btn resources-pagination-next"
                onClick={handleNextPage}
                type="button"
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
              {hasActiveFilters
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
                  onEdit={handleEditResource}
                  onDelete={handleDeleteResource}
                />
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Modals */}
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
