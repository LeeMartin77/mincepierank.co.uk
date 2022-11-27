import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getAllMakerPies,
  getAllPieRankingSummaries,
  getMincePieMakers,
  Maker,
} from "../system/storage";
import { Breadcrumbs, Divider } from "@mui/material";
import Link from "next/link";
import { PieList } from "../components/pieList/pieList";

export const getServerSideProps = async () => {
  const pies = (await getAllMakerPies()).unwrapOr([]);
  const makers = (await getMincePieMakers()).unwrapOr([]);
  const rankingSummaries = (await getAllPieRankingSummaries()).unwrapOr([]);
  return {
    props: {
      makers,
      pies,
      rankingSummaries,
    },
  };
};

function AllMincePies({
  makers,
  pies,
  rankingSummaries,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const makerMap: { [key: string]: Maker } = {};
  makers.forEach((mkr) => (makerMap[mkr.id] = mkr));
  return (
    <>
      <Head>
        <title>Mince Pie Rank :: All Pies</title>
        <meta property="og:image" content="/logo-social.png" />
      </Head>
      <main>
        <Breadcrumbs aria-label="breadcrumb">
          <Link color="inherit" href="/">
            Home
          </Link>
        </Breadcrumbs>
        <h1>All Pies</h1>
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        <PieList
          pies={pies}
          rankings={rankingSummaries}
          addMetaDescription
          metaPrefix={`All the pies in our list, ranked by you. `}
        />
      </main>
    </>
  );
}

export default AllMincePies;
