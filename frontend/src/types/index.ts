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