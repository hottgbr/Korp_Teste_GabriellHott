import { DatePipe } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { FeedbackMessage } from '../../components/feedback-message/feedback-message';
import { StatusBadge } from '../../components/status-badge/status-badge';
import { InvoiceStatus } from '../../enums/invoice-status';
import { Invoice } from '../../models/invoice';
import { Product } from '../../models/product';
import { InvoiceService } from '../../services/invoice.service';
import { ProductService } from '../../services/product.service';

@Component({
  selector: 'app-invoice-details',
  imports: [
    RouterLink,
    DatePipe,
    FeedbackMessage,
    StatusBadge,
  ],
  templateUrl: './invoice-details.html',
  styleUrl: './invoice-details.css',
})
export class InvoiceDetails implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly invoiceService = inject(InvoiceService);
  private readonly productService = inject(ProductService);

  readonly invoiceStatus = InvoiceStatus;

  invoice = signal<Invoice | null>(null);
  productsById = signal<Record<number, Product>>({});

  isLoading = signal(false);
  isClosing = signal(false);

  errorMessage = signal('');
  successMessage = signal('');
  productReferenceMessage = signal('');

  ngOnInit(): void {
    this.loadInvoice();
    this.loadProductReferences();
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

  private loadProductReferences(): void {
    this.productService
      .list()
      .subscribe({
        next: (products: Product[]) => {
          const lookup = products.reduce<Record<number, Product>>(
            (result, product) => {
              result[product.id] = product;
              return result;
            },
            {},
          );

          this.productsById.set(lookup);
        },

        error: () => {
          this.productReferenceMessage.set(
            'Não foi possível carregar os detalhes dos produtos. Os identificadores serão exibidos no lugar.',
          );
        },
      });
  }

  getProduct(productId: number): Product | undefined {
    return this.productsById()[productId];
  }

  closeInvoice(): void {
    const currentInvoice = this.invoice();

    if (
      !currentInvoice ||
      currentInvoice.status === InvoiceStatus.Closed
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

          setTimeout(() => {
            this.printInvoice();
          });
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

  printInvoice(): void {
    window.print();
  }
}