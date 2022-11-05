import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getAllMakerPies, getMincePieMakers, Maker } from "../system/storage";
import Link from "next/link";
import { Divider, Grid } from "@mui/material";
import { PieSummaryLink } from "../components/pieSummaryLink";

export const getServerSideProps = async () => {
  const pies = (await getAllMakerPies()).unwrapOr([]);
  const makers = (await getMincePieMakers()).unwrapOr([]);
  return {
    props: {
      makers,
      pies,
    },
  };
};

function AllMincePies({
  makers,
  pies,
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
        <Divider />
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
