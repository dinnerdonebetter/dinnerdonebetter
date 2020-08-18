import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class Report {
  id: number;
  reportType: string;
  concern: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.reportType = "";
    this.concern = "";
    this.createdOn = 0;
  }

static areEqual = function(
  r1: Report,
  r2: Report,
): boolean {
    return (
      r1.id === r2.id &&
      r1.reportType === r2.reportType &&
      r1.concern === r2.concern &&
      r1.archivedOn === r2.archivedOn
    );
  }
}

export const fakeReportFactory = Factory.Sync.makeFactory<Report> ({
  reportType: Factory.Sync.each(() =>  faker.random.word()),
  concern: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
