export class Pagination {
    public page: number;
    public limit: number;
    public filteredCount: number;
    public totalCount: number;

    constructor(
        page: number = 0,
        limit: number = 0,
        filteredCount: number = 0,
        totalCount: number = 0,
    ) {
        this.page = page
        this.limit = limit
        this.filteredCount = filteredCount
        this.totalCount = totalCount
    }
}
