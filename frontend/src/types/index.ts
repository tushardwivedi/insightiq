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
  type: 'superset' | 'postgres' | 'file_upload';
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
  username?: string;
  password?: string;
  bearer_token?: string;
}

export interface ConnectorTestResult {
  success: boolean;
  message: string;
  response_time?: number;
  available_datasets?: string[];
  error?: string;
}

export interface UploadedFile {
  id: string;
  user_id: string;
  filename: string;
  original_filename: string;
  file_size: number;
  file_type: 'csv' | 'excel';
  mime_type: string;
  storage_path: string;
  table_name: string;
  schema_json: {
    columns: Array<{
      name: string;
      type: string;
      nullable: boolean;
    }>;
  };
  row_count: number;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  error_message?: string;
  source: 'local' | 'google_drive';
  drive_file_id?: string;
  created_at: string;
  updated_at: string;
  processed_at?: string;
  // Ingestion tracking fields
  ingestion_status: 'pending' | 'in_progress' | 'completed' | 'failed';
  ingestion_progress: number;
  vector_count: number;
  ingestion_started_at?: string;
  ingestion_completed_at?: string;
  // Data quality metrics
  data_quality_score?: number;
  missing_data_percent?: number;
  duplicate_rows?: number;
  rows_with_missing?: number;
}