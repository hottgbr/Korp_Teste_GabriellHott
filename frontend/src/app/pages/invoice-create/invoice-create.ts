import { HttpErrorResponse } from '@angular/common/http';
import {
  Component,
  computed,
  inject,
  OnInit,
  signal,
} from '@angular/core';
import {
  FormArray,
  FormBuilder,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { FeedbackMessage } from '../../components/feedback-message/feedback-message';
import { Product } from '../../models/product';
import { InvoiceService } from '../../services/invoice.service';
import { ProductService } from '../../services/product.service';

@Component({
  selector: 'app-invoice-create',
  imports: [
    ReactiveFormsModule,
    RouterLink,
    FeedbackMessage,
  ],
  templateUrl: './invoice-create.html',
  styleUrl: './invoice-create.css',
})
export class InvoiceCreate implements OnInit {
  private readonly formBuilder = inject(FormBuilder);
  private readonly invoiceService = inject(InvoiceService);
  private readonly productService = inject(ProductService);
  private readonly router = inject(Router);

  products = signal<Product[]>([]);

  availableProducts = computed(() =>
    this.products().filter(
      (product) => product.stock > 0,
    ),
  );

  isLoadingProducts = signal(false);
  isSubmitting = signal(false);

  errorMessage = signal('');

  form = this.formBuilder.nonNullable.group({
    items: this.formBuilder.array([
      this.createItem(),
    ]),
  });

  get items(): FormArray {
    return this.form.controls.items;
  }

  ngOnInit(): void {
    this.loadProducts();
  }

  private createItem() {
    return this.formBuilder.nonNullable.group({
      productId: [
        0,
        [
          Validators.required,
          Validators.min(1),
        ],
      ],

      quantity: [
        1,
        [
          Validators.required,
          Validators.min(1),
        ],
      ],
    });
  }

  addItem(): void {
    this.items.push(
      this.createItem(),
    );
  }

  removeItem(index: number): void {
    if (this.items.length === 1) {
      return;
    }

    this.items.removeAt(index);
  }

  private loadProducts(): void {
    this.isLoadingProducts.set(true);
    this.errorMessage.set('');

    this.productService
      .list()
      .pipe(
        finalize(() => {
          this.isLoadingProducts.set(false);
        }),
      )
      .subscribe({
        next: (products: Product[]) => {
          this.products.set(products);
        },

        error: () => {
          this.errorMessage.set(
            'Não foi possível carregar os produtos.',
          );
        },
      });
  }

  private getProductStock(
    productId: number,
  ): number {
    const product = this.products().find(
      (item) => item.id === productId,
    );

    return product?.stock ?? 0;
  }

  private hasInsufficientStock(): boolean {
    const input = this.form.getRawValue();

    return input.items.some((item) => {
      const currentStock =
        this.getProductStock(
          item.productId,
        );

      return item.quantity > currentStock;
    });
  }

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const input = this.form.getRawValue();

    const productIds = input.items.map(
      (item) => item.productId,
    );

    const hasDuplicateProducts =
      new Set(productIds).size !==
      productIds.length;

    if (hasDuplicateProducts) {
      this.errorMessage.set(
        'O mesmo produto não pode ser adicionado mais de uma vez.',
      );
      return;
    }

    if (this.hasInsufficientStock()) {
      this.errorMessage.set(
        'A quantidade informada não pode ser maior que o saldo disponível do produto.',
      );
      return;
    }

    this.isSubmitting.set(true);
    this.errorMessage.set('');

    this.invoiceService
      .create(input)
      .pipe(
        finalize(() => {
          this.isSubmitting.set(false);
        }),
      )
      .subscribe({
        next: () => {
          this.router.navigate([
            '/invoices',
          ]);
        },

        error: (
          error: HttpErrorResponse,
        ) => {
          if (error.status === 0) {
            this.errorMessage.set(
              'Não foi possível conectar ao serviço de faturamento.',
            );
            return;
          }

          this.errorMessage.set(
            'Não foi possível criar a nota fiscal.',
          );
        },
      });
  }
}