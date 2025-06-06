import {
  Navigate,
  Outlet,
  RouterProvider,
  createBrowserRouter,
} from "react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { Resources } from "./pages/Resources";
import { ErrorFallback } from "./pages/ErrorFallback";

import { ReactQueryFlashProvider } from "./components/Flash";

import "./App.scss";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    element: (
      <div className="app">
        <ReactQueryFlashProvider>
          <Outlet />
        </ReactQueryFlashProvider>
      </div>
    ),
    errorElement: <ErrorFallback />,
    hydrateFallbackElement: <div>Loading...</div>,
    children: [
      {
        path: "/",
        element: <Navigate replace to="/resources" />,
      },
      {
        path: "/resources",
        element: <Resources />,
      },
    ],
  },
]);

export const Router = () => <RouterProvider router={router} />;

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router />
    </QueryClientProvider>
  );
}

export default App;
