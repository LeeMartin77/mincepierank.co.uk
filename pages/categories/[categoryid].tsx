import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getPieRankingSummariesByIds,
  getPiesWithCategory,
  MakerPie,
  PieRankingSummary,
} from "../../system/storage";
import Link from "next/link";
import { Breadcrumbs, Divider, Typography } from "@mui/material";
import { PieList } from "../../components/pieList/pieList";
import { ppCategory } from "../../components/formatCategory";

export const getServerSideProps = async ({
  params: { categoryid },
}: {
  params: { categoryid: string };
}) => {
  const pies = (await getPiesWithCategory(categoryid)).unwrapOr(
    [] as MakerPie[]
  );
  if (!pies || pies.length === 0) {
    return {
      notFound: true,
    };
  }
  const makerPies = pies.reduce((prev, pie) => {
    if (prev[pie.makerid] === undefined) {
      prev[pie.makerid] = [pie.id];
    } else {
      prev[pie.makerid].push(pie.id);
    }

    return prev;
  }, {} as { [key: string]: string[] });

  const rankingSummaries = await (
    await Promise.all(
      Object.entries(makerPies).map(([makerid, pieids]) =>
        getPieRankingSummariesByIds(makerid, pieids)
      )
    )
  ).reduce((prev, res) => {
    return [...prev, ...res.unwrapOr([] as PieRankingSummary[])];
  }, [] as PieRankingSummary[]);
  return {
    props: {
      categoryid,
      pies,
      rankingSummaries,
    },
  };
};

function Category({
  categoryid,
  pies,
  rankingSummaries,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>{`Mince Pie Rank :: ${ppCategory(categoryid)}`}</title>
        <meta property="og:image" content="/logo-social.png" />
      </Head>
      <main>
        <Breadcrumbs aria-label="breadcrumb">
          <Link color="inherit" href="/">
            Home
          </Link>
          <Typography color="text.primary">{ppCategory(categoryid)}</Typography>
        </Breadcrumbs>
        <h1>{ppCategory(categoryid)} Pies</h1>
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        <PieList
          pies={pies}
          rankings={rankingSummaries}
          addMetaDescription
          metaPrefix={`The ${ppCategory(categoryid)} pies ranked by you. `}
          lockedCategory={categoryid}
        />
      </main>
    </>
  );
}

export default Category;
