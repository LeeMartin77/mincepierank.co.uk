import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import { getMincePieMakers, Maker } from "../system/storage";
import Link from "next/link";
import { Button, Card, CardActions, CardMedia, Grid } from "@mui/material";

export const getServerSideProps = async () => {
  const data = await getMincePieMakers();
  return {
    props: {
      makers: data.unwrapOr([]),
    },
  };
};

function BrandCard(maker: Maker) {
  return (
    <Grid key={maker.id} item xs={6} sm={4} md={3}>
      <Card>
        <Link href={`/brands/${maker.id}`}>
          <CardMedia
            style={{
              objectFit: "contain",
            }}
            component="img"
            height="150"
            image={maker.logo}
            alt={`${maker.name} Logo`}
          />
        </Link>
        <CardActions>
          <Button
            LinkComponent={Link}
            href={`/brands/${maker.id}`}
            style={{ width: "100%", textAlign: "center" }}
          >
            {maker.name}
          </Button>
        </CardActions>
      </Card>
    </Grid>
  );
}

function Home({
  makers,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>Mince Pie Rank</title>
        <meta
          name="description"
          content="Rankings of UK mince pies for 2022."
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <Grid container spacing={2}>
          {makers.map(BrandCard)}
        </Grid>
      </main>
    </>
  );
}

export default Home;
