import React, { useEffect, useState } from "react";
import { Book, Video } from "lucide-react";
import { type Resource } from "../../../types";
import { resourcesApi } from "../../../services/resources";
import "./ContentDisplay.scss";

export const DefaultView: React.FC = () => {
  const [featuredResources, setFeaturedResources] = useState<Resource[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadFeaturedResources = async () => {
      try {
        const resources = await resourcesApi.getAll();
        setFeaturedResources(resources.data);
      } catch (error) {
        console.error("Failed to load featured resources:", error);
      } finally {
        setIsLoading(false);
      }
    };

    loadFeaturedResources();
  }, []);

  if (isLoading) {
    return (
      <div className="default-view">
        <div className="loading-spinner">
          <div className="spinner"></div>
          <p>Loading featured resources...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="default-view">
      <div className="welcome-section">
        <h1>Welcome to the Learning Hub</h1>
        <p>
          Explore our comprehensive collection of resources to help you make the most of our Migration App. Browse by
          topic, search for specific content, or check out our featured resources below.
        </p>
      </div>

      <div className="featured-resources">
        <h2>Resources</h2>

        <div className="featured-grid">
          {featuredResources.map((resource) => (
            <div
              key={resource.id}
              className="featured-card"
            >
              <div className="featured-icon">
                {resource.type === "video" ? <Video size={24} /> : <Book size={24} />}
              </div>
              <div className="featured-content">
                <h3>{resource.title}</h3>
                <p>{resource.description}</p>
                <div className="featured-tags">
                  {resource.tags.map((tag, index) => (
                    <span
                      key={index}
                      className="tag"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default DefaultView;
