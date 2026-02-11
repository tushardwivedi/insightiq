import {
  AnalyticsResponse,
  VoiceResponse,
  HealthCheck,
  TextQuery,
  SQLQuery,
  DataConnector,
  ConnectorTestResult,
  ApiResponse,
  UploadedFile
} from '@/types';

const API_BASE = '/api';

class ApiClient {
  private getAuthToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE}${endpoint}`;
    const token = this.getAuthToken();

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    // Add auth token if available
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    // Create AbortController for timeout - 3 minutes for LLM processing
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 180000); // 3 minutes

    try {
      const response = await fetch(url, {
        headers,
        credentials: 'same-origin',
        signal: controller.signal,
        ...options,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        // Handle 401 Unauthorized
        if (response.status === 401) {
          // Clear token and redirect to login
          if (typeof window !== 'undefined') {
            localStorage.removeItem('auth_token');
            window.location.href = '/login';
          }
        }

        const error = await response.text();
        const sanitizedError = error.replace(/<[^>]*>/g, '');
        throw new Error(`API Error: ${response.status} - ${sanitizedError}`);
      }

      if (response.status === 204 || response.headers.get('content-length') === '0') {
        return undefined as T;
      }

      return response.json();
    } catch (error) {
      clearTimeout(timeoutId);
      if (error instanceof Error && error.name === 'AbortError') {
        throw new Error('Request timeout - the operation took too long. Please try again.');
      }
      throw error;
    }
  }

  async healthCheck(): Promise<HealthCheck> {
    return this.request<HealthCheck>('/health');
  }

  async textQuery(query: TextQuery): Promise<AnalyticsResponse> {
    return this.request<AnalyticsResponse>('/query', {
      method: 'POST',
      body: JSON.stringify(query),
    });
  }

  async voiceQuery(audioFile: File): Promise<VoiceResponse> {
    const formData = new FormData();
    formData.append('audio', audioFile);

    const response = await fetch(`${API_BASE}/voice`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`Voice API Error: ${response.status} - ${error}`);
    }

    return response.json();
  }

  async sqlQuery(query: SQLQuery): Promise<AnalyticsResponse> {
    return this.request<AnalyticsResponse>('/sql', {
      method: 'POST',
      body: JSON.stringify(query),
    });
  }

  // Connector Management Methods
  async getConnectors(): Promise<ApiResponse<DataConnector[]>> {
    return this.request<ApiResponse<DataConnector[]>>('/connectors');
  }

  async createConnector(connector: Omit<DataConnector, 'id' | 'status' | 'created_at' | 'updated_at' | 'last_tested'>): Promise<DataConnector> {
    return this.request<DataConnector>('/connectors', {
      method: 'POST',
      body: JSON.stringify(connector),
    });
  }

  async updateConnector(id: string, connector: Partial<Omit<DataConnector, 'id' | 'created_at' | 'updated_at'>>): Promise<DataConnector> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<DataConnector>(`/connectors/${sanitizedId}`, {
      method: 'PUT',
      body: JSON.stringify(connector),
    });
  }

  async deleteConnector(id: string): Promise<void> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<void>(`/connectors/${sanitizedId}`, {
      method: 'DELETE',
    });
  }

  async testConnector(id: string): Promise<ConnectorTestResult> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<ConnectorTestResult>(`/connectors/${sanitizedId}/test`, {
      method: 'POST',
    });
  }

  async testConnectorConfig(config: { type: string; config: any }): Promise<ConnectorTestResult> {
    return this.request<ConnectorTestResult>('/connectors/test-config', {
      method: 'POST',
      body: JSON.stringify(config),
    });
  }

  async getConnectorData(id: string, query?: string): Promise<AnalyticsResponse> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    const params = query ? `?query=${encodeURIComponent(query)}` : '';
    return this.request<AnalyticsResponse>(`/connectors/${sanitizedId}/data${params}`);
  }

  // File Upload Methods
  async uploadFile(formData: FormData): Promise<UploadedFile> {
    const url = `${API_BASE}/files`;
    const token = this.getAuthToken();

    const headers: Record<string, string> = {};

    // Add auth token if available
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: formData,
      credentials: 'include',
    });

    if (!response.ok) {
      if (response.status === 401) {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('auth_token');
          window.location.href = '/login';
        }
      }

      const error = await response.text();
      const sanitizedError = error.replace(/<[^>]*>/g, '');
      throw new Error(`Upload Error: ${response.status} - ${sanitizedError}`);
    }

    return response.json();
  }

  async getUploadedFiles(): Promise<UploadedFile[]> {
    return this.request<UploadedFile[]>('/files');
  }

  async getUploadedFile(id: string): Promise<UploadedFile> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<UploadedFile>(`/files/${sanitizedId}`);
  }

  async deleteUploadedFile(id: string): Promise<{ success: boolean }> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<{ success: boolean }>(`/files/${sanitizedId}`, {
      method: 'DELETE',
    });
  }

  async getQueryHistory(limit?: number, offset?: number): Promise<{ data: any[]; limit: number; offset: number; count: number }> {
    const params = new URLSearchParams();
    if (limit) params.set('limit', limit.toString());
    if (offset) params.set('offset', offset.toString());
    return this.request<{ data: any[]; limit: number; offset: number; count: number }>(`/query-history?${params.toString()}`);
  }

  async getQueryHistoryStats(): Promise<any> {
    return this.request<any>('/query-history/stats');
  }

  async deleteQueryHistory(id: string): Promise<{ success: boolean; message: string }> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<{ success: boolean; message: string }>(`/query-history?id=${sanitizedId}`, {
      method: 'DELETE',
    });
  }

  async getSavedQueries(): Promise<any[]> {
    return this.request<any[]>('/saved-queries');
  }

  async saveQuery(query: { name: string; sql: string; description?: string }): Promise<any> {
    return this.request<any>('/saved-queries', {
      method: 'POST',
      body: JSON.stringify(query),
    });
  }

  async deleteSavedQuery(id: string): Promise<void> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<void>(`/saved-queries/${sanitizedId}`, {
      method: 'DELETE',
    });
  }

  async getAnalytics(): Promise<any> {
    return this.request<any>('/analytics');
  }

  async getReports(): Promise<any[]> {
    return this.request<any[]>('/reports');
  }

  async getReport(id: string): Promise<any> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<any>(`/reports/${sanitizedId}`);
  }

  async createReport(report: { name: string; description?: string; query: string }): Promise<any> {
    return this.request<any>('/reports', {
      method: 'POST',
      body: JSON.stringify(report),
    });
  }

  async deleteReport(id: string): Promise<void> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<void>(`/reports/${sanitizedId}`, {
      method: 'DELETE',
    });
  }

  async getTeamMembers(): Promise<any[]> {
    return this.request<any[]>('/team');
  }

  async inviteTeamMember(email: string, role: string): Promise<any> {
    return this.request<any>('/team/invite', {
      method: 'POST',
      body: JSON.stringify({ email, role }),
    });
  }

  async removeTeamMember(id: string): Promise<void> {
    const sanitizedId = id.replace(/[^a-zA-Z0-9-_]/g, '');
    return this.request<void>(`/team/${sanitizedId}`, {
      method: 'DELETE',
    });
  }

  async getUserSettings(): Promise<any> {
    return this.request<any>('/settings');
  }

  async updateUserSettings(settings: any): Promise<any> {
    return this.request<any>('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }
}

export const apiClient = new ApiClient();