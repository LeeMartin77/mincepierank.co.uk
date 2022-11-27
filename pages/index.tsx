import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getAllMakerPies,
  getAllPieRankingSummaries,
  getMincePieMakers,
  Maker,
} from "../system/storage";
import Link from "next/link";
import {
  Button,
  Card,
  CardActions,
  CardMedia,
  Divider,
  Grid,
} from "@mui/material";
import { mapPiesAndRankings } from "../components/mapPiesAndRankings";
import { PieSummaryLink } from "../components/pieList/pieList";
import { descriptionSummary } from "../components/descriptionSummary";
import { ppCategory } from "../components/formatCategory";

export const getServerSideProps = async () => {
  const data = (await getMincePieMakers()).unwrapOr([]);
  const pies = (await getAllMakerPies()).unwrapOr([]);
  const rankingSummaries = (await getAllPieRankingSummaries()).unwrapOr([]);
  const { mappedRankings, mappedPies, rankingOrder, unrankedPies } =
    mapPiesAndRankings(pies, rankingSummaries);
  const topPieId = rankingOrder.shift();
  return {
    props: {
      makers: data,
      topPie: topPieId ? mappedPies[topPieId] : undefined,
      topPieRanking: topPieId ? mappedRankings[topPieId] : undefined,
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
  topPieRanking,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const topPieMaker = makers.find((x) => x.id === topPie?.makerid);
  const description =
    topPie && topPieMaker && topPieRanking
      ? ` ${descriptionSummary(topPie, topPieRanking, topPieMaker)}`
      : ``;
  const mainCategories = [
    "classic",
    "premium",
    "alternative",
    "vegan",
    "mini",
    "gluten-free",
  ];
  return (
    <>
      <Head>
        <title>Mince Pie Rank</title>
        <meta
          name="description"
          content={`Rankings of UK mince pies for 2022.${description}`}
        />
        <meta property="og:image" content="/logo.svg" />
      </Head>
      <main>
        {topPie && topPieRanking && (
          <>
            <h1>Current Top Pie</h1>
            <PieSummaryLink isTop pie={topPie} ranking={topPieRanking} />
          </>
        )}
        <Button
          LinkComponent={Link}
          href={`/all-mince-pies`}
          style={{ width: "100%", textAlign: "center", marginTop: "1em" }}
        >
          All Mince Pies
        </Button>
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        <h2>Categories</h2>
        <Grid container spacing={2}>
          {mainCategories.map((categoryid) => (
            <Grid key={categoryid} item xs={6} sm={4}>
              <Button
                LinkComponent={Link}
                href={`/categories/${categoryid}`}
                style={{ width: "100%", textAlign: "center" }}
              >
                {ppCategory(categoryid)}
              </Button>
            </Grid>
          ))}
        </Grid>
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        <h2>Brands</h2>
        <Grid container spacing={2}>
          {makers.map(BrandCard)}
        </Grid>
      </main>
    </>
  );
}

export default Home;
