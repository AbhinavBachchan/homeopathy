import { Component, computed, inject, signal } from "@angular/core";
import { NgIf } from "@angular/common";
import { NavigationEnd, Router, RouterOutlet } from "@angular/router";
import { filter } from "rxjs";
import { HeaderComponent } from "./shared/header/header.component";
import { FooterComponent } from "./shared/footer/footer.component";

@Component({
  selector: "app-root",
  standalone: true,
  imports: [RouterOutlet, NgIf, HeaderComponent, FooterComponent],
  templateUrl: "./app.component.html",
  styleUrl: "./app.component.css",
})
export class AppComponent {
}
