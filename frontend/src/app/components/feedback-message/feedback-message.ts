import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-feedback-message',
  templateUrl: './feedback-message.html',
  styleUrl: './feedback-message.css',
})
export class FeedbackMessage {
  @Input({ required: true }) message = '';
  @Input() type: 'success' | 'error' | 'info' = 'info';
}