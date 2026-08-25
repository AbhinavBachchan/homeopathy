import { TestBed } from "@angular/core/testing";
import { provideHttpClient } from "@angular/common/http";
import { provideRouter } from "@angular/router";
import { ProductListComponent } from "./product-list.component";

describe("ProductListComponent", () => {
  it("creates", async () => {
    await TestBed.configureTestingModule({
      imports: [ProductListComponent],
      providers: [provideHttpClient(), provideRouter([])],
    }).compileComponents();
    expect(
      TestBed.createComponent(ProductListComponent).componentInstance,
    ).toBeTruthy();
  });
});
