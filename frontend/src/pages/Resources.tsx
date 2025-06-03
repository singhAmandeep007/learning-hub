import { useState } from "react";
import { FilterTabs } from "../components/FilterTabs";
import { SearchBar } from "../components/SearchBar";
import { ResourceList } from "./components/ResourceList";
import { ContentDisplay } from "./components/ContentDisplay";
import { useResources } from "../services/resources/hooks";
import { useTags } from "../services/tags/hooks";
import type { Resource, ResourceType } from "../types";
import "./Resources.scss";

const Resources = () => {
  const [activeResource, setActiveResource] = useState<Resource | null>(null);
  const [activeType, setActiveType] = useState<ResourceType | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [activeTag, setActiveTag] = useState<string | null>(null);

  const { data: resourcesData, isLoading: isLoadingResources } = useResources({
    search: searchQuery,
    type: activeType ?? undefined,
    tags: activeTag ? [activeTag] : undefined,
  });

  const { data: tags } = useTags();

  const handleSearch = (query: string) => {
    setSearchQuery(query);
  };

  const handleTypeFilter = (type: ResourceType | null) => {
    setActiveType(type);
  };

  const handleResourceSelect = (resource: Resource) => {
    setActiveResource(resource);
  };

  return (
    <div className="learning-hub">
      <header className="learning-hub-header">
        <h1>Learning Hub</h1>
        <p>Find tutorials, guides, and resources to help you get the most out of the Migration App.</p>
      </header>

      <div className="learning-hub-container">
        <aside className="sidebar">
          <div className="sidebar-content">
            <SearchBar onSearch={handleSearch} />

            <FilterTabs
              activeFilter={activeType}
              onFilterChange={handleTypeFilter}
            />

            <div className="resources-container">
              <h3 className="resources-heading">Resources</h3>
              <ResourceList
                resources={resourcesData?.data ?? []}
                activeResourceId={activeResource?.id || null}
                onResourceSelect={handleResourceSelect}
                isLoading={isLoadingResources}
              />
            </div>
          </div>
        </aside>

        <main className="main-content">
          <ContentDisplay activeResource={activeResource} />
        </main>
      </div>
    </div>
  );
};

export default Resources;