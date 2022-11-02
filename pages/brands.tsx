import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMakers } from "../system/storage";

export const getServerSideProps = async () => {
  const data = await getMincePieMakers();
  return {
    props: {
      // Probably a smarter way of doing this
      makers: data.isOk() ? data._unsafeUnwrap() : [],
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
              <li key={mkr.id}>{mkr.name}</li>
            ))}
          </ul>
        )}
      </main>
    </>
  );
}

export default Brands;
