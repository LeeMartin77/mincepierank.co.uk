import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMakers } from "../system/storage";
import Link from "next/link";

export const getServerSideProps = async () => {
  const data = await getMincePieMakers();
  return {
    props: {
      makers: data.unwrapOr([]),
    },
  };
};

function Brands({
  makers,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>Mince Pie Rank :: Brands</title>
        <meta
          name="description"
          content="The overarching brands of pie maker we have in our database"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        {makers.length > 0 && (
          <ul>
            {makers.map((mkr) => (
              <li key={mkr.id}>
                <Link href={`/brands/${mkr.id}`}>{mkr.name}</Link>
              </li>
            ))}
          </ul>
        )}
      </main>
    </>
  );
}

export default Brands;
