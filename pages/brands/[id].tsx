import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMaker } from "../../system/storage";

export const getServerSideProps = async ({
  params: { id },
}: {
  params: { id: string };
}) => {
  const data = await getMincePieMaker(id);
  return {
    props: {
      maker: data.unwrapOr(undefined),
    },
  };
};

function Brands({
  maker,
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
      </main>
    </>
  );
}

export default Brands;
