import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMaker, getPiesByMaker } from "../../system/storage";
import Link from "next/link";

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
        {maker && <>{maker.name}</>}
        {!maker && <>Maker Not Found</>}
        <ul>
          {pies.map((pie) => (
            <li key={pie.id}>
              <Link href={`/brands/${pie.makerid}/${pie.id}`}>
                {pie.displayname}
              </Link>
            </li>
          ))}
        </ul>
      </main>
    </>
  );
}

export default Brands;
