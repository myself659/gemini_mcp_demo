const API_BASE_URL = "http://localhost:8081/api";

interface ApiResponse<T> {
  data?: T;
  error?: string;
}

export async function post<T>(path: string, body: any): Promise<ApiResponse<T>> {
  try {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    const data = await response.json();

    if (!response.ok) {
      return { error: data.error || "An unknown error occurred" };
    }

    return { data };
  } catch (error: any) {
    return { error: error.message || "Network error" };
  }
}

export async function get<T>(path: string, token?: string): Promise<ApiResponse<T>> {
  try {
    const headers: HeadersInit = {};
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${path}`, {
      method: "GET",
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      return { error: data.error || "An unknown error occurred" };
    }

    return { data };
  } catch (error: any) {
    return { error: error.message || "Network error" };
  }
}
