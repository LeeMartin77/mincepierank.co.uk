import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getAllMakerPies,
  getAllPieRankingSummaries,
  getMincePieMakers,
} from "../system/storage";
import Link from "next/link";
import { Button, Divider, Grid } from "@mui/material";
import { mapPiesAndRankings } from "../components/mapPiesAndRankings";
import { descriptionSummary } from "../components/descriptionSummary";
import { ppCategory } from "../components/formatCategory";
import { PieSummaryLink } from "../components/pieList/pieSummaryLink";
import { BrandCard } from "../components/brandCard";

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
        <meta property="og:image" content="/logo-social.png" />
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
          style={{
            width: "100%",
            textAlign: "center",
            marginTop: "1em",
            height: "4em",
          }}
        >
          Browse All Mince Pies
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
          {makers.map((mkr) => (
            <BrandCard key={mkr.id} maker={mkr} />
          ))}
        </Grid>
      </main>
    </>
  );
}

export default Home;
