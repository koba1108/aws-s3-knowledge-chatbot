import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ChatRequest, ChatResponse } from '../models/chat.model';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private apiUrl = environment.apiUrl;
  private sessionId?: string;

  constructor(private http: HttpClient) {}

  sendMessage(message: string, knowledgeBaseId?: string): Observable<ChatResponse> {
    const request: ChatRequest = {
      message,
      session_id: this.sessionId,
      knowledge_base_id: knowledgeBaseId
    };

    return this.http.post<ChatResponse>(`${this.apiUrl}/api/chat`, request);
  }

  setSessionId(sessionId: string): void {
    this.sessionId = sessionId;
  }

  getSessionId(): string | undefined {
    return this.sessionId;
  }

  clearSession(): void {
    this.sessionId = undefined;
  }

  healthCheck(): Observable<any> {
    return this.http.get(`${this.apiUrl}/api/health`);
  }
}
