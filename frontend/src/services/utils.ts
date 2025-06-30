import { type Product } from "../types";

/**
 * Extracts the product parameter from the current URL path
 * Expected URL format: /:product/resources
 */
export function getProductFromUrl(): Product {
  const path = window.location.pathname;
  const segments = path.split("/").filter(Boolean);

  return segments[0] as Product;
}
