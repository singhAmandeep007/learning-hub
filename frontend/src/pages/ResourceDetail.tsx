import { useParams } from "react-router";
import { useResource } from "../services/resources/hooks";
import "./ResourceDetail.scss";

const ResourceDetail = () => {
  const { id } = useParams<{ id: string }>();
  const { data: resource, isLoading } = useResource({
    id: id!,
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!resource) {
    return <div>Resource not found</div>;
  }

  return (
    <div className="resourceDetail">
      <header className="header"></header>
      <h1 className="title">{resource.title}</h1>
      <div className="tags">
        {resource.tags.map((tag) => (
          <span
            key={tag}
            className="tag"
          >
            {tag}
          </span>
        ))}
      </div>
      <div className="description">{resource.description}</div>
      <div className="media">
        {resource.type === "video" ? (
          <video
            src={resource.url}
            controls
          />
        ) : resource.type === "pdf" ? (
          <embed
            src={resource.url}
            type="application/pdf"
            width="100%"
            height="600px"
          />
        ) : (
          <iframe
            src={resource.url}
            width="100%"
            height="600px"
            title={resource.title}
          />
        )}
      </div>
    </div>
  );
};

export default ResourceDetail;
