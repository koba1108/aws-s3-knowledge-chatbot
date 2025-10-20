import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ChatService } from '../../services/chat.service';
import { ChatMessage } from '../../models/chat.model';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.css']
})
export class ChatComponent implements OnInit {
  messages: ChatMessage[] = [];
  userMessage = '';
  isLoading = false;
  isConnected = false;

  constructor(private chatService: ChatService) {}

  ngOnInit(): void {
    this.checkConnection();
  }

  checkConnection(): void {
    this.chatService.healthCheck().subscribe({
      next: () => {
        this.isConnected = true;
      },
      error: () => {
        this.isConnected = false;
      }
    });
  }

  sendMessage(): void {
    if (!this.userMessage.trim() || this.isLoading) {
      return;
    }

    const userMsg: ChatMessage = {
      role: 'user',
      content: this.userMessage,
      timestamp: new Date()
    };

    this.messages.push(userMsg);
    const messageToSend = this.userMessage;
    this.userMessage = '';
    this.isLoading = true;

    this.chatService.sendMessage(messageToSend).subscribe({
      next: (response) => {
        const assistantMsg: ChatMessage = {
          role: 'assistant',
          content: response.response,
          timestamp: new Date(),
          sources: response.sources
        };
        this.messages.push(assistantMsg);
        this.chatService.setSessionId(response.session_id);
        this.isLoading = false;
      },
      error: (error) => {
        const errorMsg: ChatMessage = {
          role: 'assistant',
          content: `Error: ${error.error?.error || error.message || 'Failed to get response'}`,
          timestamp: new Date()
        };
        this.messages.push(errorMsg);
        this.isLoading = false;
      }
    });
  }

  clearChat(): void {
    this.messages = [];
    this.chatService.clearSession();
  }

  onKeyPress(event: KeyboardEvent): void {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      this.sendMessage();
    }
  }
}
