import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMaker, getPiesByMaker } from "../../system/storage";

export const getServerSideProps = async ({
  params: { id },
}: {
  params: { id: string };
}) => {
  const maker = (await getMincePieMaker(id)).unwrapOr(undefined);
  const pies = (await getPiesByMaker(id)).unwrapOr([]);
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
        <title>Mince Pie Rank :: {maker ? maker.name : "Not Found"}</title>
        <meta
          name="description"
          content="The overarching brands of pie maker we have in our database"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        {maker && <>{maker.name}</>}
        {!maker && <>Maker Not Found</>}
        <ul>
          {pies.map((pie) => (
            <li key={pie.id}>{pie.displayname}</li>
          ))}
        </ul>
      </main>
    </>
  );
}

export default Brands;
