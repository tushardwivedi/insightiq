import {
  AnalyticsResponse,
  VoiceResponse,
  HealthCheck,
  TextQuery,
  SQLQuery,
  DataConnector,
  ConnectorTestResult,
  ApiResponse
} from '@/types';

const API_BASE = '/api';

class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE}${endpoint}`;

    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      credentials: 'same-origin',
      ...options,
    });

    if (!response.ok) {
      const error = await response.text();
      const sanitizedError = error.replace(/<[^>]*>/g, '');
      throw new Error(`API Error: ${response.status} - ${sanitizedError}`);
    }

    return response.json();
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
}

export const apiClient = new ApiClient();