import { HttpErrorResponse } from '@angular/common/http';
import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { FeedbackMessage } from '../../components/feedback-message/feedback-message';
import { Invoice } from '../../models/invoice';
import { InvoiceService } from '../../services/invoice.service';

@Component({
  selector: 'app-invoice-details',
  imports: [
    RouterLink,
    FeedbackMessage,
  ],
  templateUrl: './invoice-details.html',
  styleUrl: './invoice-details.css',
})
export class InvoiceDetails implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly invoiceService = inject(InvoiceService);

  invoice = signal<Invoice | null>(null);

  isLoading = signal(false);
  isClosing = signal(false);

  errorMessage = signal('');
  successMessage = signal('');

  ngOnInit(): void {
    this.loadInvoice();
  }

  private loadInvoice(): void {
    const id = Number(
      this.route.snapshot.paramMap.get('id'),
    );

    if (!id) {
      this.errorMessage.set(
        'Identificador da nota fiscal inválido.',
      );
      return;
    }

    this.isLoading.set(true);
    this.errorMessage.set('');
    this.successMessage.set('');

    this.invoiceService
      .findById(id)
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (invoice: Invoice) => {
          this.invoice.set(invoice);
        },

        error: (error: HttpErrorResponse) => {
          if (error.status === 404) {
            this.errorMessage.set(
              'Nota fiscal não encontrada.',
            );
            return;
          }

          this.errorMessage.set(
            'Não foi possível carregar a nota fiscal.',
          );
        },
      });
  }

  closeInvoice(): void {
    const currentInvoice = this.invoice();

    if (
      !currentInvoice ||
      currentInvoice.status === 'CLOSED'
    ) {
      return;
    }

    this.isClosing.set(true);
    this.errorMessage.set('');
    this.successMessage.set('');

    this.invoiceService
      .close(currentInvoice.id)
      .pipe(
        finalize(() => {
          this.isClosing.set(false);
        }),
      )
      .subscribe({
        next: (invoice: Invoice) => {
          this.invoice.set(invoice);
          this.successMessage.set(
            'Nota fiscal fechada com sucesso. O estoque foi atualizado.',
          );
        },

        error: (error: HttpErrorResponse) => {
          switch (error.status) {
            case 404:
              this.errorMessage.set(
                'Produto ou nota fiscal não encontrada.',
              );
              break;

            case 409:
              this.errorMessage.set(
                'Não há estoque suficiente para fechar esta nota fiscal.',
              );
              break;

            case 503:
              this.errorMessage.set(
                'O serviço de estoque está indisponível no momento.',
              );
              break;

            default:
              this.errorMessage.set(
                'Não foi possível fechar a nota fiscal.',
              );
          }
        },
      });
  }
}