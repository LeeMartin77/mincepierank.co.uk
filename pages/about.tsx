import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMakers } from "../system/storage";

function About() {
  return (
    <>
      <Head>
        <title>Mince Pie Rank :: About</title>
        <meta name="description" content="All about Mince Pie Rank" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <div>Pithy text to come</div>
      </main>
    </>
  );
}

export default About;
