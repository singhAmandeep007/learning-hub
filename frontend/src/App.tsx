import { Navigate, Outlet, RouterProvider, createBrowserRouter, useParams } from "react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

import { Ban } from "lucide-react";

import { Resources } from "./pages/Resources";
import { ErrorFallback } from "./pages/ErrorFallback";
import { NotFound } from "./pages/NotFound";

import { ReactQueryFlashProvider } from "./components/Flash";
import { VALID_PRODUCTS, type Product, DEFAULT_PRODUCT } from "./types";

import styles from "./App.module.scss";

const queryClient = new QueryClient();

const productResourceLoader = ({
  params,
}: {
  params: {
    product?: string;
  };
}) => {
  const { product } = params;

  // Check if the 'product' parameter is in our list of valid product names
  if (!VALID_PRODUCTS.includes(product as Product)) {
    throw new Response("Invalid Product Resource", { status: 404 });
  }
  return null;
};

const InvalidProductPage = () => {
  const { product } = useParams<{ product?: string }>();

  return (
    <div className={styles.invalidProduct}>
      <Ban size={48} />
      <h1>Invalid Product</h1>
      <p>The requested product "{product}" is not valid.</p>
    </div>
  );
};

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
            to={`/${DEFAULT_PRODUCT}/resources`}
          />
        ),
      },
      {
        path: "/:product/resources",
        element: <Resources />,
        loader: productResourceLoader,
        errorElement: <InvalidProductPage />,
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
