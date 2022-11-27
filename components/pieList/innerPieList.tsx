import { Divider, Grid } from "@mui/material";
import { MakerPie } from "../../system/storage";
import { PieListRanking, PieSummaryLink } from "./pieSummaryLink";

export default function InnerPieList({
  rankingOrder,
  mappedPies,
  mappedRankings,
  unrankedPies,
}: {
  mappedRankings: {
    [key: string]: PieListRanking | undefined;
  };
  mappedPies: { [key: string]: MakerPie };
  rankingOrder: string[];
  unrankedPies: Set<string>;
}) {
  return (
    <>
      <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />

      <Grid container spacing={2}>
        {rankingOrder.map((uniqid) => (
          <Grid item key={uniqid} xs={12} md={6} lg={4}>
            <PieSummaryLink
              pie={mappedPies[uniqid]}
              ranking={mappedRankings[uniqid]}
            />
          </Grid>
        ))}
        {Array.from(unrankedPies).map((uniqid) => {
          return (
            <Grid item key={uniqid} xs={12} md={6} lg={4}>
              <PieSummaryLink pie={mappedPies[uniqid]} />
            </Grid>
          );
        })}
      </Grid>
    </>
  );
}
