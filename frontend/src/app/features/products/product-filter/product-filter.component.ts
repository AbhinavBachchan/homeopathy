import { Component, EventEmitter, Output, signal } from '@angular/core';
import { CommonModule } from '@angular/common';

export interface FilterOption {
  id: string;
  label: string;
  value: string;
  count?: number;
}

export interface FilterCategory {
  id: string;
  name: string;
  description?: string;
  options: FilterOption[];
}

// Mock array of objects for categories and their checkbox options
// In future, these will be fetched directly from the backend API
export const MOCK_FILTER_DATA: FilterCategory[] = [
  {
    id: 'homeopathy',
    name: 'Homeopathy',
    description: 'Formulation categories',
    options: [
      { id: 'medicines', label: 'Medicines', value: 'Medicines', count: 355 },
      { id: 'dilutions', label: 'Dilutions', value: 'Dilutions', count: 109 },
      { id: 'mother_tinctures', label: 'Mother Tinctures', value: 'Mother Tinctures', count: 98 },
      { id: 'triturations', label: 'Triturations', value: 'Triturations', count: 0 },
      { id: 'biochemic', label: 'Biochemic', value: 'Biochemic', count: 0 },
      { id: 'bio_combination', label: 'Bio Combination', value: 'Bio Combination', count: 0 },
      { id: 'cosmetics', label: 'Cosmetics', value: 'Cosmetics', count: 6 },
    ],
  },
  {
    id: 'potency',
    name: 'Potency',
    description: 'Select remedy potencies',
    options: [
      { id: 'q', label: 'Mother Tincture (Q)', value: 'Q', count: 98 },
      { id: '30c', label: '30C', value: '30C', count: 109 },
      { id: '200c', label: '200C', value: '200C', count: 85 },
      { id: '1m', label: '1M', value: '1M', count: 42 },
      { id: '6ch', label: '6 CH', value: '6 CH', count: 30 },
      { id: '3x', label: '3X', value: '3X', count: 15 },
    ],
  },
  {
    id: 'brand',
    name: 'Brand / Manufacturer',
    description: 'Filter by certified manufacturers',
    options: [
      { id: 'sbl', label: 'SBL Homeopathy', value: 'SBL', count: 140 },
      { id: 'dr_reckeweg', label: 'Dr. Reckeweg (Germany)', value: 'Dr. Reckeweg', count: 120 },
      { id: 'schwabe', label: 'Schwabe India (WSG)', value: 'Schwabe India', count: 85 },
      { id: 'baksons', label: "Bakson's Homeopathy", value: "Bakson's", count: 45 },
      { id: 'adel', label: 'Adel Pekana', value: 'Adel Pekana', count: 32 },
    ],
  },
];

export interface ActiveFilters {
  [categoryId: string]: string[];
}

@Component({
  selector: 'app-product-filter',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './product-filter.component.html',
  styleUrl: './product-filter.component.css',
})
export class ProductFilterComponent {
  @Output() filterChange = new EventEmitter<ActiveFilters>();

  // Available categories array of objects
  categories = signal<FilterCategory[]>(MOCK_FILTER_DATA);

  // Store selected checkbox values by categoryId: { [categoryId: string]: string[] }
  selectedFilterMap = signal<ActiveFilters>({});

  /**
   * Toggle a checkbox for a specific category and option value.
   */
  toggleOption(categoryId: string, optionValue: string, event: Event): void {
    const isChecked = (event.target as HTMLInputElement).checked;
    const currentMap = { ...this.selectedFilterMap() };
    const currentValues = currentMap[categoryId] ? [...currentMap[categoryId]] : [];

    if (isChecked) {
      if (!currentValues.includes(optionValue)) {
        currentValues.push(optionValue);
      }
    } else {
      const idx = currentValues.indexOf(optionValue);
      if (idx > -1) {
        currentValues.splice(idx, 1);
      }
    }

    if (currentValues.length > 0) {
      currentMap[categoryId] = currentValues;
    } else {
      delete currentMap[categoryId];
    }

    this.selectedFilterMap.set(currentMap);
    this.filterChange.emit(currentMap);
  }

  isOptionSelected(categoryId: string, optionValue: string): boolean {
    const currentMap = this.selectedFilterMap();
    return !!currentMap[categoryId]?.includes(optionValue);
  }

  getSelectedCountForCategory(categoryId: string): number {
    return this.selectedFilterMap()[categoryId]?.length || 0;
  }

  getTotalActiveFiltersCount(): number {
    const currentMap = this.selectedFilterMap();
    return Object.values(currentMap).reduce((acc, curr) => acc + curr.length, 0);
  }

  /**
   * Remove single selected filter item
   */
  removeFilter(categoryId: string, optionValue: string): void {
    const currentMap = { ...this.selectedFilterMap() };
    if (currentMap[categoryId]) {
      currentMap[categoryId] = currentMap[categoryId].filter((v) => v !== optionValue);
      if (currentMap[categoryId].length === 0) {
        delete currentMap[categoryId];
      }
      this.selectedFilterMap.set(currentMap);
      this.filterChange.emit(currentMap);
    }
  }

  /**
   * Clear all selected filters
   */
  clearAllFilters(): void {
    this.selectedFilterMap.set({});
    this.filterChange.emit({});
  }

  /**
   * Clear filters for a single category
   */
  clearCategoryFilters(categoryId: string): void {
    const currentMap = { ...this.selectedFilterMap() };
    delete currentMap[categoryId];
    this.selectedFilterMap.set(currentMap);
    this.filterChange.emit(currentMap);
  }

  /**
   * Helper to get active filter tags for badge display
   */
  getActiveFilterTags(): Array<{ categoryId: string; categoryName: string; value: string; label: string }> {
    const tags: Array<{ categoryId: string; categoryName: string; value: string; label: string }> = [];
    const map = this.selectedFilterMap();
    const categoriesList = this.categories();

    for (const [catId, values] of Object.entries(map)) {
      const cat = categoriesList.find((c) => c.id === catId);
      const catName = cat ? cat.name : catId;

      for (const val of values) {
        const opt = cat?.options.find((o) => o.value === val);
        tags.push({
          categoryId: catId,
          categoryName: catName,
          value: val,
          label: opt ? opt.label : val,
        });
      }
    }
    return tags;
  }
}
