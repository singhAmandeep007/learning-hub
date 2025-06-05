export class HttpClient {
  private baseURL: string;
  private defaultHeaders: Record<string, string>;

  constructor(baseURL: string = "/api") {
    this.baseURL = baseURL.replace(/\/$/, ""); // Remove trailing slash
    this.defaultHeaders = {
      "Content-Type": "application/json",
    };
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;

    const config: RequestInit = {
      headers: {
        ...this.defaultHeaders,
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      // Handle empty responses (like 204 No Content)
      if (response.status === 204 || response.headers.get("content-length") === "0") {
        return {} as T;
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("An unexpected error occurred");
    }
  }

  async get<T>(endpoint: string, params?: Record<string, string | number | boolean | string[]>): Promise<T> {
    const processedParams = params
      ? Object.entries(params).reduce(
          (acc, [key, value]) => {
            acc[key] = Array.isArray(value) ? value.join(",") : value;
            return acc;
          },
          {} as Record<string, string | number | boolean>
        )
      : undefined;

    const url = processedParams
      ? `${endpoint}?${new URLSearchParams(processedParams as Record<string, string>)}`
      : endpoint;
    return this.request<T>(url, { method: "GET" });
  }

  async post<T>(endpoint: string, data?: Record<string, unknown>): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async postFormData<T>(endpoint: string, options: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      headers: {},
      ...options,
    });
  }

  async patch<T>(endpoint: string, data: Record<string, unknown>): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  }

  async patchFormData<T>(endpoint: string, options: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PATCH",
      headers: {},
      ...options,
    });
  }

  async put<T>(endpoint: string, data: Record<string, unknown>): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" });
  }

  // Method to set authorization header
  setAuthToken(token: string) {
    this.defaultHeaders["Authorization"] = `Bearer ${token}`;
  }

  // Method to remove authorization header
  removeAuthToken() {
    delete this.defaultHeaders["Authorization"];
  }

  // Method to add custom headers
  setHeader(key: string, value: string) {
    this.defaultHeaders[key] = value;
  }
}

export const httpClient = new HttpClient();
