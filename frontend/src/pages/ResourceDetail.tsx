import React from "react";
import { useParams } from "react-router";

import { type Resource } from "../types";

import "./ResourceDetail.scss";
import { resourcesApi } from "../services/resources";

const ResourceDetail = () => {
  const { id } = useParams<{ id: string }>();

  const [resource, setResource] = React.useState<Resource | null>(null);
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    if (id) {
      resourcesApi.getById({ id }).then((data) => {
        setResource(data);
        setLoading(false);
      });
    }
  }, [id]);

  return (
    <div className="resourceDetail">
      <header className="header"></header>

      <h1 className="title">resource.title</h1>

      <div className="tags">Resources</div>

      <div className="description">resource.description</div>

      <div className="media">Media</div>
    </div>
  );
};

export default ResourceDetail;
