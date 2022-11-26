import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getAllMakerPies,
  getAllPieRankingSummaries,
  getMincePieMakers,
  getUserPieRankings,
  Maker,
} from "../system/storage";
import { Breadcrumbs, Button, Divider } from "@mui/material";
import Link from "next/link";
import { PieList } from "../components/pieList/pieList";
import { getSession, signIn } from "next-auth/react";

export const getServerSideProps = async (ctx: any) => {
  const session = await getSession(ctx);
  const user = session?.user?.email;
  if (!user) {
    return { props: { loggedIn: false } };
  }
  const rankings = (await getUserPieRankings(user)).unwrapOr([]);
  const pies = (await getAllMakerPies()).unwrapOr([]);
  return {
    props: {
      loggedIn: true,
      pies: pies.filter(
        (x) =>
          rankings.findIndex(
            (y) => y.makerid === x.makerid && y.pieid === x.id
          ) > -1
      ),
      rankings,
    },
  };
};

function AllMincePies({
  loggedIn,
  pies,
  rankings,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>Mince Pie Rank :: My Rankings</title>
        <meta name="description" content="All my pie rankings" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <Breadcrumbs aria-label="breadcrumb">
          <Link color="inherit" href="/">
            Home
          </Link>
        </Breadcrumbs>
        <h1>My Pie Rankings</h1>
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        {!loggedIn && (
          <Button
            style={{ width: "100%", textAlign: "center" }}
            onClick={() =>
              signIn(
                process.env.NODE_ENV === "development"
                  ? "Development"
                  : "google"
              )
            }
          >
            Sign in with{" "}
            {process.env.NODE_ENV === "development" ? "Development" : "Google"}
          </Button>
        )}
        {loggedIn && pies && rankings && rankings.length > 0 && (
          <PieList pies={pies} rankings={rankings} />
        )}
        {loggedIn && (!rankings || rankings.length === 0) && (
          <div style={{ width: "100%", textAlign: "center" }}>
            You&apos;ve not ranked any pies yet
          </div>
        )}
      </main>
    </>
  );
}

export default AllMincePies;
