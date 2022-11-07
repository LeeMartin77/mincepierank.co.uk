import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getAllMakerPies, getAllPieRankingSummaries, getMincePieMakers, Maker } from "../system/storage";
import { Button, Card, CardActions, CardHeader, CardMedia, Divider, Grid } from "@mui/material";
import { PieSummaryLink } from "../components/pieSummaryLink";
import { findTopPie } from "../components/findTopPie";
import Link from "next/link";

export const getServerSideProps = async () => {
  const pies = (await getAllMakerPies()).unwrapOr([]);
  const makers = (await getMincePieMakers()).unwrapOr([]);
  const rankingSummaries = (await getAllPieRankingSummaries()).unwrapOr([]);
  const [topPie, topPieRanking] = findTopPie(pies, rankingSummaries) ?? [null, null];
  if (topPie && topPieRanking) {
    return {
      props: {
        makers,
        pies,
        topPie,
        topPieRanking
      },
    };
  }
  return {
    props: {
      makers,
      pies
    },
  };
};

function AllMincePies({
  makers,
  pies,
  topPie,
  topPieRanking
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const makerMap: { [key: string]: Maker } = {};
  makers.forEach((mkr) => (makerMap[mkr.id] = mkr));
  return (
    <>
      <Head>
        <title>Mince Pie Rank :: All Pies</title>
        <meta
          name="description"
          content="A big, interactive list of all the pies in our database"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <h1>All Pies</h1>
        {topPie && topPieRanking && <Card>
          <CardHeader title="The Top Pie"/>
          <Link href={`/brands/${topPie.makerid}/${topPie.id}`}>
            <CardMedia
              component="img"
              height="200"
              image={topPie.image_file}
              alt={`${topPie.displayname}`}
            />
          </Link>
          <CardActions>
            <Button
              LinkComponent={Link}
              href={`/brands/${topPie.makerid}/${topPie.id}`}
              style={{ width: "100%", textAlign: "center" }}
            >
              {topPie.displayname}
            </Button>
          </CardActions>
        </Card>}
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        <Grid container spacing={2}>
          {pies.map((pie) => {
            return (
              <Grid key={pie.id + pie.makerid} item xs={6} sm={4} md={3}>
                <PieSummaryLink pie={pie} />
              </Grid>
            );
          })}
        </Grid>
      </main>
    </>
  );
}

export default AllMincePies;
