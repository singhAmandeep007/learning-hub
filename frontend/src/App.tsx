import { Navigate, Outlet, RouterProvider, createBrowserRouter } from "react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

import { Resources } from "./pages/Resources";
import { ErrorFallback } from "./pages/ErrorFallback";
import { NotFound } from "./pages/NotFound";

import { ReactQueryFlashProvider } from "./components/Flash";

import styles from "./App.module.scss";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    element: (
      <div className={styles.app}>
        <ReactQueryFlashProvider>
          <Outlet />
        </ReactQueryFlashProvider>
      </div>
    ),
    errorElement: <ErrorFallback />,
    hydrateFallbackElement: <div className={styles.loading}>Loading...</div>,
    children: [
      {
        path: "/",
        element: (
          <Navigate
            replace
            to="/resources"
          />
        ),
      },
      {
        path: "/resources",
        element: <Resources />,
      },
      {
        path: "*",
        element: <NotFound />,
      },
    ],
  },
]);

export const Router = () => <RouterProvider router={router} />;

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router />
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}

export default App;
