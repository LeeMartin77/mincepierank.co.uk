import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getAllRankingsForPie,
  getMincePieMaker,
  getPieByMakerAndId,
} from "../../../system/storage";

export const getServerSideProps = async ({
  params: { brandid, pieid },
}: {
  params: { brandid: string; pieid: string };
}) => {
  const maker = (await getMincePieMaker(brandid)).unwrapOr(undefined);
  const pie = (await getPieByMakerAndId(brandid, pieid)).unwrapOr(undefined);
  const rankings = (await getAllRankingsForPie(brandid, pieid)).unwrapOr([]);
  return {
    props: {
      maker,
      pie,
      rankings,
    },
  };
};

function Brands({
  maker,
  pie,
  rankings,
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
        {maker && <>{maker.name}</>}
        {!maker && <>Maker Not Found</>}
        {pie && <span>{pie.displayname}</span>}
        <ul>
          {rankings.map((rank) => {
            return <li key={rank.userid}>{rank.notes}</li>;
          })}
        </ul>
      </main>
    </>
  );
}

export default Brands;
