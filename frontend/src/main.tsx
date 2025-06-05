import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";

async function enableMocking() {
  if (process.env.NODE_ENV !== "development") {
    return;
  }
return Promise.resolve();
  // const { worker } = await import("./mocks/browser");

  // // `worker.start()` returns a Promise that resolves
  // // once the Service Worker is up and ready to intercept requests.
  // return worker.start();
}

enableMocking().then(() => {
  createRoot(document.getElementById("root")!).render(
    <StrictMode>
      <App />
    </StrictMode>
  );
});
