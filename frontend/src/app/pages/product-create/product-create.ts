import { HttpErrorResponse } from '@angular/common/http';
import { Component, inject, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { finalize } from 'rxjs';


import { ProductService } from '../../services/product.service';

@Component({
  selector: 'app-product-create',
  imports: [ ReactiveFormsModule, RouterLink, ],
  templateUrl: './product-create.html',
  styleUrl: './product-create.css',
})
export class ProductCreate {
  private readonly formBuilder = inject(FormBuilder);
  private readonly productService = inject(ProductService);
  private readonly router = inject(Router);

  isSubmitting = signal(false);
  errorMessage = signal('');

  form = this.formBuilder.nonNullable.group({
    code: [
      '',
      [
        Validators.required,
      ],
    ],
    description: [
      '',
      [
        Validators.required,
      ],
    ],
    stock: [
      0,
      [
        Validators.required,
        Validators.min(0),
      ],
    ],
  });

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.isSubmitting.set(true);
    this.errorMessage.set('');

    this.productService
      .create(this.form.getRawValue())
      .pipe(
        finalize(() => {
          this.isSubmitting.set(false);
        }),
      )
      .subscribe({
        next: () => {
          this.router.navigate(['/products']);
        },

        error: (error: HttpErrorResponse) => {
          if (error.status === 409) {
            this.errorMessage.set(
              'Já existe um produto cadastrado com este código.',
            );
            return;
          }

          if (error.status === 0) {
            this.errorMessage.set(
              'Não foi possível conectar ao serviço de estoque.',
            );
            return;
          }

          this.errorMessage.set(
            'Não foi possível cadastrar o produto. Tente novamente.',
          );
        },
      });
  }
}
