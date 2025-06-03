import React from "react";
import { type Resource } from "../../../types";

import "./ContentDisplay.scss";

interface VideoPlayerProps {
  resource: Resource;
}

export const VideoPlayer: React.FC<VideoPlayerProps> = ({ resource }) => {
  return (
    <div className="video-player">
      <div className="video-header">
        <h2>{resource.title}</h2>
        <div className="video-meta">
          {resource.tags.length > 0 && (
            <div className="video-tags">
              {resource.tags.map((tag, index) => (
                <span
                  key={index}
                  className="tag"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="video-container">
        {/* This is a placeholder for the actual video embed */}
        <div className="video-placeholder">
          <div className="placeholder-content">
            <div className="play-icon">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="48"
                height="48"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <polygon points="5 3 19 12 5 21 5 3"></polygon>
              </svg>
            </div>
            <p>Video Player: {resource.title}</p>
            <p className="placeholder-note">
              In the actual implementation, this would be a real video player embedded from Zoom.
            </p>
          </div>
        </div>
      </div>

      <div className="video-description">
        <h3>Description</h3>
        <p>{resource.description}</p>
      </div>
    </div>
  );
};
