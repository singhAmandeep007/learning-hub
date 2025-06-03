import React from "react";
import { type Resource } from "../../../types";

import { VideoPlayer } from "./VideoPlayer";
import { ArticleViewer } from "./ArticleViewer";
import { DefaultView } from "./DefaultView";

import "./ContentDisplay.scss";

interface ContentDisplayProps {
  activeResource: Resource | null;
}

export const ContentDisplay: React.FC<ContentDisplayProps> = ({ activeResource }) => {
  if (!activeResource) {
    return <DefaultView />;
  }

  return (
    <div className="content-display">
      {activeResource.type === "video" ? (
        <VideoPlayer resource={activeResource} />
      ) : (
        <ArticleViewer resource={activeResource} />
      )}
    </div>
  );
};
