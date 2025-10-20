export interface ChatMessage {
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  sources?: Source[];
}

export interface Source {
  content: string;
  location: {
    uri?: string;
  };
}

export interface ChatRequest {
  message: string;
  session_id?: string;
  knowledge_base_id?: string;
}

export interface ChatResponse {
  response: string;
  session_id: string;
  sources?: Source[];
  error?: string;
}
