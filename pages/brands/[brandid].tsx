import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getMakerPieRankingSummaries,
  getMincePieMaker,
  getPiesByMaker,
  MakerPie,
  PieRankingSummary,
} from "../../system/storage";
import Link from "next/link";
import {
  Breadcrumbs,
  Button,
  Card,
  CardActions,
  CardHeader,
  CardMedia,
  Divider,
  Grid,
  Typography,
} from "@mui/material";
import { PieList } from "../../components/pieList/pieList";

export const getServerSideProps = async ({
  params: { brandid },
}: {
  params: { brandid: string };
}) => {
  const maker = (await getMincePieMaker(brandid)).unwrapOr(undefined);
  if (!maker) {
    return {
      notFound: true,
    };
  }

  const pies = (await getPiesByMaker(brandid)).unwrapOr([] as MakerPie[]);
  const rankingSummaries = (
    await getMakerPieRankingSummaries(brandid)
  ).unwrapOr([] as PieRankingSummary[]);
  return {
    props: {
      maker,
      pies,
      rankingSummaries,
    },
  };
};

function Brands({
  maker,
  pies,
  rankingSummaries,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>{`Mince Pie Rank :: ${maker.name}`}</title>
      </Head>
      <main>
        <Breadcrumbs aria-label="breadcrumb">
          <Link color="inherit" href="/">
            Home
          </Link>
          <Typography color="text.primary">{maker.name}</Typography>
        </Breadcrumbs>
        <h1>{maker.name}</h1>
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        <PieList
          pies={pies}
          rankings={rankingSummaries}
          addMetaDescription
          metaPrefix={`The pies of ${maker.name} ranked by you. `}
        />
      </main>
    </>
  );
}

export default Brands;
