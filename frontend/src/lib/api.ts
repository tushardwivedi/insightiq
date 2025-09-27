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
      ...options,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`API Error: ${response.status} - ${error}`);
    }

    return response.json();
  }

  async healthCheck(): Promise<HealthCheck> {
    return this.request<HealthCheck>('/health');
  }

  async textQuery(query: TextQuery): Promise<AnalyticsResponse> {
    return this.request<AnalyticsResponse>('/direct-query', {
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
    return this.request<DataConnector>(`/connectors/${id}`, {
      method: 'PUT',
      body: JSON.stringify(connector),
    });
  }

  async deleteConnector(id: string): Promise<void> {
    return this.request<void>(`/connectors/${id}`, {
      method: 'DELETE',
    });
  }

  async testConnector(id: string): Promise<ConnectorTestResult> {
    return this.request<ConnectorTestResult>(`/connectors/${id}/test`, {
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
    const params = query ? `?query=${encodeURIComponent(query)}` : '';
    return this.request<AnalyticsResponse>(`/connectors/${id}/data${params}`);
  }
}

export const apiClient = new ApiClient();