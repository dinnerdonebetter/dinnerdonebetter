// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPagination {
  page: number;
  totalCount: number;
  filteredCount: number;
  limit: number;
}

export class Pagination implements IPagination {
  page: number;
  totalCount: number;
  filteredCount: number;
  limit: number;
  constructor(input: Partial<Pagination> = {}) {
    this.page = input.page = 0;
    this.totalCount = input.totalCount = 0;
    this.filteredCount = input.filteredCount = 0;
    this.limit = input.limit = 0;
  }
}
