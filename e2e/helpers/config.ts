export const product = process.env.E2E_PRODUCT ?? "ecomm";
export const apiBaseURL = process.env.E2E_API_BASE_URL ?? "http://localhost:8000";

export const resourcesCollectionPath = `/api/v1/${product}/resources`;

export const resourceByIdPath = (id: string) => `${resourcesCollectionPath}/${id}`;
