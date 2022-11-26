import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { 
  getAllMakerPies,
  getAllPieRankingSummaries,
  getMincePieMakers, Maker } from "../system/storage";
import Link from "next/link";
import { Button, Card, CardActions, CardMedia, Divider, Grid } from "@mui/material";
import { mapPiesAndRankings } from "../components/mapPiesAndRankings";
import { PieSummaryLink } from "../components/pieList/pieList";

export const getServerSideProps = async () => {
  const data = (await getMincePieMakers()).unwrapOr([]);
  const pies = (await getAllMakerPies()).unwrapOr([]);
  const rankingSummaries = (await getAllPieRankingSummaries()).unwrapOr([]);
  const { mappedRankings, mappedPies, rankingOrder, unrankedPies } = mapPiesAndRankings(pies, rankingSummaries)
  const topPieId = rankingOrder.shift();
  return {
    props: {
      makers: data,
      topPie: topPieId ? mappedPies[topPieId] : undefined,
      topPieRanking: topPieId ? mappedRankings[topPieId] : undefined
    },
  };
};

function BrandCard(maker: Maker) {
  return (
    <Grid key={maker.id} item xs={6} sm={4} md={3}>
      <Card>
        <Link href={`/brands/${maker.id}`}>
          <CardMedia
            style={{
              objectFit: "contain",
            }}
            component="img"
            height="150"
            image={maker.logo}
            alt={`${maker.name} Logo`}
          />
        </Link>
        <CardActions>
          <Button
            LinkComponent={Link}
            href={`/brands/${maker.id}`}
            style={{ width: "100%", textAlign: "center" }}
          >
            {maker.name}
          </Button>
        </CardActions>
      </Card>
    </Grid>
  );
}

function Home({
  makers,
  topPie,
  topPieRanking
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>Mince Pie Rank</title>
        <meta
          name="description"
          content={`Rankings of UK mince pies for 2022.${topPie ? ` The current leader is ${topPie.displayname} from ${makers.find(x => x.id === topPie.makerid)?.name} with ${topPieRanking?.average.toFixed(1)} stars across ${topPieRanking?.count} votes.` : ``}`}
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        {topPie && topPieRanking && (
          <>
            <h1>Current Top Pie</h1>
            <PieSummaryLink
              isTop
              pie={topPie}
              ranking={topPieRanking}
            />
          </>
        )}
        <Divider style={{ marginTop: '1em', marginBottom: '1em' }}/>
        <Grid container spacing={2}>
          {makers.map(BrandCard)}
        </Grid>
      </main>
    </>
  );
}

export default Home;
