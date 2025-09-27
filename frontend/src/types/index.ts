export interface HealthCheck {
  status: 'healthy' | 'unhealthy';
  services: {
    [key: string]: {
      status: 'healthy' | 'unhealthy';
      message?: string;
    };
  };
  timestamp?: string;
}

export interface AnalyticsResponse {
  data: any[];
  query: string;
  insights?: string;
  timestamp?: string;
  process_time?: string;
  task_id?: string;
  status?: string;
}

export interface VoiceResponse {
  transcript: string;
  response: AnalyticsResponse;
  task_id: string;
  process_time: string;
  status: string;
}

export interface TextQuery {
  query: string;
}

export interface SQLQuery {
  sql: string;
  question?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface DataConnector {
  id: string;
  name: string;
  type: 'superset' | 'postgres' | 'mysql' | 'mongodb' | 'api';
  status: 'connected' | 'disconnected' | 'error' | 'testing';
  config: ConnectorConfig;
  created_at: string;
  updated_at: string;
  last_tested?: string;
}

export interface ConnectorConfig {
  url: string;
  username?: string;
  password?: string;
  api_key?: string;
  bearer_token?: string;
  database?: string;
  additional_params?: Record<string, any>;
}

export interface SupersetConnectorConfig extends ConnectorConfig {
  username: string;
  password: string;
  bearer_token?: string;
}

export interface ConnectorTestResult {
  success: boolean;
  message: string;
  response_time?: number;
  available_datasets?: string[];
  error?: string;
}