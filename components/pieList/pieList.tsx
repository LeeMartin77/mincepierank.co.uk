import { MakerPie, PieRankingSummary } from "../../system/storage";
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardMedia,
  Grid,
  Link,
  Rating,
  Typography,
} from "@mui/material";
import Image from "next/image";

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
            }}
          >
            <Image
              src={pie.image_file}
              alt={pie.displayname}
              height={200}
              width={200}
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
  rankings.forEach((ranking) => {
    mappedRankings[ranking.makerid + "-" + ranking.pieid] = ranking;
  });
  return (
    <Grid container spacing={2}>
      {pies.map((pie) => (
        <Grid item key={pie.makerid + pie.id} xs={12} md={6} lg={4}>
          <PieSummaryLink
            pie={pie}
            ranking={mappedRankings[pie.makerid + "-" + pie.id]}
          />
        </Grid>
      ))}
    </Grid>
  );
}
