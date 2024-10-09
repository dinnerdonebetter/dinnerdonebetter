// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IPagination {
   filteredCount: number;
 limit: number;
 page: number;
 totalCount: number;

}

export class Pagination implements IPagination {
   filteredCount: number;
 limit: number;
 page: number;
 totalCount: number;
constructor(input: Partial<Pagination> = {}) {
	 this.filteredCount = input.filteredCount = 0;
 this.limit = input.limit = 0;
 this.page = input.page = 0;
 this.totalCount = input.totalCount = 0;
}
}