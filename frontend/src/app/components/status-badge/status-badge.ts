import { Component, Input } from '@angular/core';

import { InvoiceStatus } from '../../enums/invoice-status';

@Component({
  selector: 'app-status-badge',
  templateUrl: './status-badge.html',
  styleUrl: './status-badge.css',
})
export class StatusBadge {
  @Input({ required: true })
  status!: InvoiceStatus;

  readonly invoiceStatus = InvoiceStatus;
}