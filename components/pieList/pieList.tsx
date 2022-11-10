/* eslint-disable @next/next/no-img-element */
import { MakerPie, PieRankingSummary } from "../../system/storage";
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardMedia,
  Divider,
  Grid,
  Link,
  Rating,
  Typography,
} from "@mui/material";

export function PieSummaryLink({
  pie,
  ranking,
}: {
  pie: MakerPie;
  ranking?: PieRankingSummary;
}) {
  return (
    <Card data-testid={pie.id + "-link-card"}>
      <CardContent>
        <div style={{ display: "flex" }}>
          <div
            style={{
              width: "47.5%",
              height: "120px",
            }}
          >
            <img
              src={pie.image_file}
              alt={pie.displayname}
              style={{
                width: "100%",
                height: "100%",
                objectFit: "cover",
              }}
            />
          </div>
          <div
            style={{
              marginLeft: "5%",
              width: "47.5%",
            }}
          >
            <Typography component="legend">Avg. Ranking</Typography>
            <Rating name="read-only" value={ranking?.average ?? 0} readOnly />
            <div style={{ paddingLeft: "0.25em" }}>
              <Typography variant="caption">
                {ranking?.count ?? 0} Rankings
              </Typography>
            </div>
          </div>
        </div>
      </CardContent>
      <CardActions>
        <Button
          LinkComponent={Link}
          href={`/brands/${pie.makerid}/${pie.id}`}
          style={{ width: "100%", textAlign: "center" }}
        >
          {pie.displayname}
        </Button>
      </CardActions>
    </Card>
  );
}

export function PieList({
  pies,
  rankings,
}: {
  pies: MakerPie[];
  rankings: PieRankingSummary[];
}) {
  const mappedRankings: { [key: string]: PieRankingSummary | undefined } = {};
  const mappedPies: { [key: string]: MakerPie } = {};
  const rankingOrder: string[] = [];

  const unrankedPies: Set<string> = new Set();
  pies.forEach((pie) => {
    const uniqId = pie.makerid + "-" + pie.id;
    mappedPies[uniqId] = pie;
    unrankedPies.add(uniqId);
  });
  rankings
    .sort((a, b) => b.average - a.average)
    .forEach((ranking) => {
      const uniqId = ranking.makerid + "-" + ranking.pieid;
      rankingOrder.push(uniqId);
      unrankedPies.delete(uniqId);
      mappedRankings[uniqId] = ranking;
    });
  const topPieId = rankingOrder.shift();
  return (
    <>
      {topPieId && mappedPies[topPieId] && (
        <>
          <h1>Top Pie</h1>
          <PieSummaryLink
            pie={mappedPies[topPieId]}
            ranking={mappedRankings[topPieId]}
          />
        </>
      )}
      {!topPieId && (
        <Typography
          variant="h4"
          style={{
            width: "100%",
            marginTop: "1em",
            marginBottom: "1em",
            textAlign: "center",
          }}
        >
          No top pie yet...
        </Typography>
      )}
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
          console.log(uniqid);
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
