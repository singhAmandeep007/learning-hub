import { Outlet, RouterProvider, createBrowserRouter } from "react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import Resources from "./pages/Resources";
import ResourceDetail from "./pages/ResourceDetail";
import AdminPanel from "./pages/AdminPanel";

import { CreateUpdateResource } from "./pages/CreateUpdateResource";

import "./App.scss";
import { ReactQueryFlashProvider } from "./components/Flash";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    element: (
      <div className="app">
        <Outlet />
      </div>
    ),
    children: [
      {
        path: "/",
        element: <Resources />,
      },
      {
        path: "/resource/:id",
        element: <ResourceDetail />,
      },
      {
        path: "/admin",
        element: <AdminPanel />,
      },
      {
        path: "/create-resource",
        element: <CreateUpdateResource />,
      },
    ],
  },
]);

export const Router = () => <RouterProvider router={router} />;

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ReactQueryFlashProvider>
        <Router />
      </ReactQueryFlashProvider>
    </QueryClientProvider>
  );
}

export default App;
