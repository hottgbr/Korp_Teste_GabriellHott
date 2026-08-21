import { Component, inject } from '@angular/core';
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

  isSubmitting = false;
  errorMessage = '';

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

    this.isSubmitting = true;
    this.errorMessage = '';

    this.productService
      .create(this.form.getRawValue())
      .pipe(
        finalize(() => {
          this.isSubmitting = false;
        }),
      )
      .subscribe({
        next: () => {
          this.router.navigate(['/products']);
        },

        error: () => {
          this.errorMessage =
            'Não foi possível cadastrar o produto.';
        },
      });
  }
}
