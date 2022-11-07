import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMakerPieRankingSummaries, getMincePieMaker, getPiesByMaker } from "../../system/storage";
import Link from "next/link";
import { Button, Card, CardActions, CardHeader, CardMedia, Divider, Grid } from "@mui/material";
import { PieSummaryLink } from "../../components/pieSummaryLink";
import { findTopPie } from "../../components/findTopPie";

export const getServerSideProps = async ({
  params: { brandid },
}: {
  params: { brandid: string };
}) => {
  const maker = (await getMincePieMaker(brandid)).unwrapOr(undefined);
  const pies = (await getPiesByMaker(brandid)).unwrapOr([]);
  const rankingSummaries = (await getMakerPieRankingSummaries(brandid)).unwrapOr([]);
  const [topPie, topPieRanking] = findTopPie(pies, rankingSummaries) ?? [];
  if (topPie && topPieRanking) {
    return {
      props: {
        maker,
        pies,
        topPie,
        topPieRanking
      },
    };
  }
  return {
    props: {
      maker,
      pies
    },
  };
};

function Brands({
  maker,
  pies,
  topPie,
  topPieRanking
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>{`Mince Pie Rank :: ${maker ? maker.name : "Not Found"}`}</title>
        <meta
          name="description"
          content="The overarching brands of pie maker we have in our database"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        {maker && <h1>{maker.name}</h1>}
        {topPie && topPieRanking && maker && <Card>
          <CardHeader title={`${maker.name} Best Pie`}/>
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

export default Brands;
