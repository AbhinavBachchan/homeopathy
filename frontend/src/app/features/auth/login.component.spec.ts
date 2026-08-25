import { TestBed } from "@angular/core/testing";
import { provideHttpClient } from "@angular/common/http";
import { provideRouter } from "@angular/router";
import { LoginComponent } from "./login.component";

describe("LoginComponent", () => {
  it("creates", async () => {
    await TestBed.configureTestingModule({
      imports: [LoginComponent],
      providers: [provideHttpClient(), provideRouter([])],
    }).compileComponents();
    expect(
      TestBed.createComponent(LoginComponent).componentInstance,
    ).toBeTruthy();
  });
});
