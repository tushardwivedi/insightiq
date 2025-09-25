import { AnalyticsResponse, VoiceResponse, HealthCheck, TextQuery, SQLQuery } from '@/types';

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
}

export const apiClient = new ApiClient();