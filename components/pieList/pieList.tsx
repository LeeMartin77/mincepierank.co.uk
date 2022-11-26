/* eslint-disable @next/next/no-img-element */
import {
  MakerPie,
  MakerPieRanking,
  PieRankingSummary,
} from "../../system/storage";
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
import { formatPrice } from "../formatPrice";
import { mapPiesAndRankings } from "../mapPiesAndRankings";
import Head from "next/head";
import { descriptionSummary } from "../descriptionSummary";

export type PieListRanking = Omit<MakerPieRanking, "userid" | "notes"> & {
  count: number | undefined;
  average: number;
};

export function PieSummaryLink({
  pie,
  ranking,
  isTop = false,
}: {
  pie: MakerPie;
  ranking?: PieListRanking;
  isTop?: boolean;
}) {
  return (
    <Card
      data-testid={pie.id + "-link-card"}
      elevation={isTop ? 8 : 1}
      style={
        !isTop
          ? {}
          : {
              color: "black",
              fontWeight: "bold",
              background:
                "linear-gradient(29deg, rgba(255,160,27,1) 0%, rgba(255,209,144,1) 60%, rgba(252,176,69,1) 100%)",
            }
      }
    >
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
                objectFit: "contain",
              }}
            />
          </div>
          <div
            style={{
              marginLeft: "5%",
              width: "47.5%",
            }}
          >
            <Typography
              component="legend"
              style={isTop ? { fontWeight: "bold" } : {}}
            >
              Avg. Ranking
            </Typography>
            <Rating
              style={
                isTop
                  ? {
                      boxShadow: "2px 2px 5px rgba(0,0,0,0.4)",
                      backgroundColor: "rgba(0,0,0,0.4)",
                    }
                  : {}
              }
              name="read-only"
              value={ranking?.average ?? 0}
              readOnly
            />
            {ranking && ranking.count != undefined && (
              <div style={{ paddingLeft: "0.25em" }}>
                <Typography
                  variant="caption"
                  style={isTop ? { fontWeight: "bold" } : {}}
                >
                  {ranking?.count ?? 0} Rankings
                </Typography>
              </div>
            )}
            <div style={{ paddingLeft: "0.25em" }}>
              <Typography
                variant="caption"
                style={isTop ? { fontWeight: "bold" } : {}}
              >
                {formatPrice(pie.pack_price_in_pence / pie.pack_count)} per pie
              </Typography>
            </div>
          </div>
        </div>
      </CardContent>
      <CardActions>
        <Button
          LinkComponent={Link}
          variant={isTop ? "outlined" : undefined}
          href={`/brands/${pie.makerid}/${pie.id}`}
          style={
            isTop
              ? {
                  width: "100%",
                  textAlign: "center",
                  color: "black",
                  borderColor: "black",
                  fontWeight: "bold",
                }
              : { width: "100%", textAlign: "center" }
          }
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
  addMetaDescription = false,
  metaPrefix = ""
}: {
  pies: MakerPie[];
  rankings: PieListRanking[];
  addMetaDescription?: boolean;
  metaPrefix?: string;
}) {
  const { mappedRankings, mappedPies, rankingOrder, unrankedPies } = mapPiesAndRankings(pies, rankings)
  const topPieId = rankingOrder.shift();
  return (
    <>
      {addMetaDescription && <Head>
        <meta
          name="description"
          content={`${metaPrefix}${topPieId ? descriptionSummary(mappedPies[topPieId]!, mappedRankings[topPieId]!) : ''}`}
        />
      </Head>}
      {topPieId && mappedPies[topPieId] && (
        <>
          <h1>Top Pie</h1>
          <PieSummaryLink
            isTop
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
