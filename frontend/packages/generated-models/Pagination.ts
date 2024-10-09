// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPagination {
  limit: number;
  page: number;
  totalCount: number;
  filteredCount: number;
}

export class Pagination implements IPagination {
  limit: number;
  page: number;
  totalCount: number;
  filteredCount: number;
  constructor(input: Partial<Pagination> = {}) {
    this.limit = input.limit = 0;
    this.page = input.page = 0;
    this.totalCount = input.totalCount = 0;
    this.filteredCount = input.filteredCount = 0;
  }
}
