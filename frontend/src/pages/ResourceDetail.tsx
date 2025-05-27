import { useParams, Link } from "react-router";
import styles from "./ResourceDetail.module.scss";

const ResourceDetail = () => {
  const { id } = useParams<{ id: string }>();

  console.log("id=", id);

  return (
    <div className={styles.resourceDetail}>
      <header className={styles.header}>
        <Link
          to="/"
          className={styles.backLink}
        >
          ← Back to Resources
        </Link>
      </header>

      <main className={styles.content}>
        <h1 className={styles.title}>resource.title</h1>

        <div className={styles.tags}>Resources</div>

        <div className={styles.description}>resource.description</div>

        <div className={styles.media}>Media</div>
      </main>
    </div>
  );
};

export default ResourceDetail;
