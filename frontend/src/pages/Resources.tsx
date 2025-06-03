import { useEffect, useState } from "react";

import { FilterTabs } from "../components/FilterTabs";
import { SearchBar } from "../components/SearchBar";

import { ResourceList } from "./components/ResourceList";

import "./Resources.scss";

import type { Resource, ResourceType, Tag } from "../types";

import { resourcesApi } from "../services/resources";
import { tagsApi } from "../services/tags";
import { Search } from "lucide-react";
import { ContentDisplay } from "./components/ContentDisplay";

const Resources = () => {
  const [resources, setResources] = useState<Resource[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [activeResource, setActiveResource] = useState<Resource | null>(null);
  const [activeType, setActiveType] = useState<ResourceType | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [activeTag, setActiveTag] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      try {
        setIsLoading(true);

        // Load topics
        const tagsData = await tagsApi.getAll();

        // Load Resources
        const resourcesData = await resourcesApi.getAll({
          search: searchQuery,
          type: activeType ? activeType : undefined,
          tags: activeTag ? [activeTag] : undefined,
        });

        setTags(tagsData);

        setResources(resourcesData.data);
      } catch (error) {
        console.error("Failed to load data:", error);
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [searchQuery, activeType, activeTag]);

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
        <aside className={`sidebar`}>
          <div className="sidebar-content">
            <SearchBar onSearch={handleSearch} />

            <FilterTabs
              activeFilter={activeType}
              onFilterChange={handleTypeFilter}
            />

            <div className="resources-container">
              <h3 className="resources-heading">Resources</h3>
              <ResourceList
                resources={resources}
                activeResourceId={activeResource?.id || null}
                onResourceSelect={handleResourceSelect}
                isLoading={isLoading}
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
