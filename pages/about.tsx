import Head from "next/head";
import { Button, Divider, Link } from "@mui/material";
import { signOut, useSession } from "next-auth/react";

function About() {
  const { status } = useSession();

  return (
    <>
      <Head>
        <title>Mince Pie Rank :: About</title>
        <meta name="description" content="All about Mince Pie Rank" />
      </Head>
      <main>
        <h1>Mince Pie Rank</h1>
        <h3>What is this?</h3>
        <p>
          This started as a fun project between me and my partner - where we
          wanted to rank the mince pies we ate over christmas. What was
          interesting was the second we mentioned the idea to other people, they
          were interested too.
        </p>
        <h3>Cookie and Privacy Policy</h3>
        <p>
          We use a simple authentication cookie once a user is logged in, and
          only store a user identifier so that users know what pies they have
          ranked. We store no other user data than this, and I&apos;ve got no
          intention of storing more.
        </p>
        <p>
          <Link href="https://github.com/LeeMartin77/mincepierank.co.uk">
            Source code is publicly available on Github.
          </Link>
        </p>
        <h4>Issues? Questions?</h4>
        <p>
          Email{" "}
          <Link href="mailto:piemaster@mincepierank.co.uk">
            piemaster@mincepierank.co.uk
          </Link>
        </p>
        {status !== "unauthenticated" && status !== "loading" && (
          <>
            <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
            <Button
              color="secondary"
              style={{ width: "100%", textAlign: "center" }}
              onClick={() => signOut()}
            >
              Log Out
            </Button>
          </>
        )}
      </main>
    </>
  );
}

export default About;
