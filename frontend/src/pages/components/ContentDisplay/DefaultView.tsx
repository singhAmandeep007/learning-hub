import React from "react";
import { Book, Video } from "lucide-react";
import { useResources } from "../../../services/resources/hooks";
import "./ContentDisplay.scss";

export const DefaultView: React.FC = () => {
  const { data: resourcesData, isLoading } = useResources();

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
          {resourcesData?.data.map((resource) => (
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