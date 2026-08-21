import { Component, inject, OnInit, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { Invoice } from '../../models/invoice';
import { InvoiceService } from '../../services/invoice.service';

@Component({
  selector: 'app-invoices',
  imports: [RouterLink],
  templateUrl: './invoices.html',
  styleUrl: './invoices.css',
})
export class Invoices implements OnInit {
  private readonly invoiceService = inject(InvoiceService);

  invoices = signal<Invoice[]>([]);
  isLoading = signal(false);
  errorMessage = signal('');

  ngOnInit(): void {
    this.loadInvoices();
  }

  loadInvoices(): void {
    this.isLoading.set(true);
    this.errorMessage.set('');

    this.invoiceService
      .list()
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (invoices: Invoice[]) => {
          this.invoices.set(invoices);
        },

        error: () => {
          this.errorMessage.set(
            'Não foi possível carregar as notas fiscais.',
          );
        },
      });
  }
}