import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMaker, getPiesByMaker } from "../../system/storage";
import Link from "next/link";
import { Divider, Grid } from "@mui/material";
import { PieSummaryLink } from "../../components/pieSummaryLink";

export const getServerSideProps = async ({
  params: { brandid },
}: {
  params: { brandid: string };
}) => {
  const maker = (await getMincePieMaker(brandid)).unwrapOr(undefined);
  const pies = (await getPiesByMaker(brandid)).unwrapOr([]);
  return {
    props: {
      maker,
      pies,
    },
  };
};

function Brands({
  maker,
  pies,
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

export default Brands;
