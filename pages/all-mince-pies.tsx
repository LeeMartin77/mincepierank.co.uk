import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getAllMakerPies, getMincePieMakers, Maker } from "../system/storage";
import Link from "next/link";

export const getServerSideProps = async () => {
  const pies = (await getAllMakerPies()).unwrapOr([]);
  const makers = (await getMincePieMakers()).unwrapOr([]);
  return {
    props: {
      makers,
      pies,
    },
  };
};

function AllMincePies({
  makers,
  pies,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const makerMap: { [key: string]: Maker } = {};
  makers.forEach((mkr) => (makerMap[mkr.id] = mkr));
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
        {pies.length > 0 && (
          <ul>
            {pies.map((pie) => {
              return (
                <li key={pie.id}>
                  <Link href={`/brands/${pie.makerid}/${pie.id}`}>
                    {pie.displayname}
                  </Link>
                </li>
              );
            })}
          </ul>
        )}
      </main>
    </>
  );
}

export default AllMincePies;
