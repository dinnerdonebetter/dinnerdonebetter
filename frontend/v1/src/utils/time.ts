import * as moment from "moment";

export function renderUnixTime(time?: number): string {
  if (time && time !== 0) {
    return moment.unix(time).format("MM/DD/YYYY");
  }
  return "never";
}
