import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getMincePieMaker,
  getPieByMakerAndId,
  getPieRankingSummary,
} from "../../../system/storage";
import {
  Card,
  CardContent,
  CardHeader,
  Divider,
  Rating,
  Typography,
} from "@mui/material";

export const getServerSideProps = async ({
  params: { brandid, pieid },
}: {
  params: { brandid: string; pieid: string };
}) => {
  const maker = (await getMincePieMaker(brandid)).unwrapOr(undefined);
  const pie = (await getPieByMakerAndId(brandid, pieid)).unwrapOr(undefined);
  const rankingSummary = (await getPieRankingSummary(brandid, pieid)).unwrapOr(
    undefined
  );
  return {
    props: {
      maker,
      pie,
      rankingSummary,
    },
  };
};

function RankingSummary({ label, value }: { label: string; value: number }) {
  return (
    <>
      <Typography component="legend">
        {`${label} (${value.toFixed(2)})`}
      </Typography>
      <Rating name="read-only" value={value} readOnly />
    </>
  );
}

function Brands({
  maker,
  pie,
  rankingSummary,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>
          {`Mince Pie Rank :: ${maker ? maker.name : "Not Found"} :: ${
            pie ? pie.displayname : "Not Found"
          }`}
        </title>
        <meta
          name="description"
          content="The overarching brands of pie maker we have in our database"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <h1>
          {maker && maker.name} :: {pie && pie.displayname}
        </h1>
        <Divider />
        {rankingSummary && (
          <Card>
            <CardHeader
              title="Summary"
              subheader={`${rankingSummary.count} Rankings`}
            />
            <CardContent>
              <RankingSummary label="Filling" value={rankingSummary.filling} />
              <RankingSummary label="Pastry" value={rankingSummary.pastry} />
              <RankingSummary label="Topping" value={rankingSummary.topping} />
              <RankingSummary label="Looks" value={rankingSummary.looks} />
              <RankingSummary label="Value" value={rankingSummary.value} />
            </CardContent>
          </Card>
        )}
      </main>
    </>
  );
}

export default Brands;
