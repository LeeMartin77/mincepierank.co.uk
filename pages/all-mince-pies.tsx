import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMakers } from "../system/storage";

export const getServerSideProps = async () => {
  const data = await getMincePieMakers();
  return {
    props: {
      makers: data.unwrapOr([]),
    },
  };
};

function AllMincePies({
  makers,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>Mince Pie Rank :: All Pies</title>
        <meta
          name="description"
          content="A big, interactive list of all the pies in our database"
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

export default AllMincePies;
