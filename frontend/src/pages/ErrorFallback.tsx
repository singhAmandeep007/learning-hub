import { isRouteErrorResponse, useRouteError } from "react-router";

import { TriangleAlert } from "lucide-react";

import "./ErrorFallback.scss";

export const ErrorFallback = () => {
  const error = useRouteError();

  let errorMessage: string;

  if (isRouteErrorResponse(error)) {
    errorMessage = `${error?.status} ${error.statusText}`;
  } else if (error instanceof Error && error.message) {
    errorMessage = error.message;
  } else if (typeof error === "string") {
    errorMessage = error;
  } else {
    console.error(error);
    errorMessage = "Please click reload to load the app again!";
  }

  return (
    <div className="error-fallback">
      <TriangleAlert size={48} className="error-fallback-icon" />
      <div>
        <h1 className="error-fallback-caption">Something went wrong!</h1>
      </div>

      <div className="error-fallback-message">
        <pre>{errorMessage}</pre>
      </div>
      <button
        className="error-fallback-button"
        onClick={() => (window.location.href = "/")}
      >
        Reload App
      </button>
    </div>
  );
};
