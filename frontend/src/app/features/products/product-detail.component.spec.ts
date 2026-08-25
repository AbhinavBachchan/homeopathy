import { TestBed } from "@angular/core/testing";
import { provideHttpClient } from "@angular/common/http";
import { ActivatedRoute } from "@angular/router";
import { ProductDetailComponent } from "./product-detail.component";

describe("ProductDetailComponent", () => {
  it("creates", async () => {
    await TestBed.configureTestingModule({
      imports: [ProductDetailComponent],
      providers: [
        provideHttpClient(),
        {
          provide: ActivatedRoute,
          useValue: { snapshot: { paramMap: new Map() } },
        },
      ],
    }).compileComponents();
    expect(
      TestBed.createComponent(ProductDetailComponent).componentInstance,
    ).toBeTruthy();
  });
});
