import { MakerPie } from "../system/storage";

export function PieSummaryInformation({ pie }: { pie: MakerPie }) {
  return <div>{pie.displayname}</div>;
}
