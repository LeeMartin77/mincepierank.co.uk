import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getMincePieMaker,
  getPieByMakerAndId,
  getPieRankingSummary,
  MakerPieRanking,
} from "../../../system/storage";
import {
  Alert,
  Breadcrumbs,
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  CardMedia,
  Divider,
  Rating,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import { signIn, signOut, useSession } from "next-auth/react";
import Image from "next/image";
import Link from "next/link";

export const getServerSideProps = async ({
  params: { brandid, pieid },
}: {
  params: { brandid: string; pieid: string };
}) => {
  const maker = (await getMincePieMaker(brandid)).unwrapOr(undefined);
  const pie = (await getPieByMakerAndId(brandid, pieid)).unwrapOr(undefined);
  if (!maker || !pie) {
    return {
      notFound: true,
    };
  }
  const rankingSummary = (await getPieRankingSummary(brandid, pieid)).unwrapOr(
    undefined
  );
  return {
    props: {
      maker,
      pie,
      rankingSummary,
    },
  };
};

function RankingSummary({
  label,
  value,
  setValue,
}: {
  label: string;
  value: number;
  setValue?: (n: number) => void;
}) {
  return (
    <>
      <Typography component="legend">
        {`${label} (${value.toFixed(2)})`}
      </Typography>
      <Rating
        name="read-only"
        value={value}
        readOnly={!setValue}
        onChange={(event, newValue) => {
          setValue && newValue && setValue(newValue);
        }}
      />
    </>
  );
}

type PieRankingDetails = Pick<
  MakerPieRanking,
  "filling" | "pastry" | "topping" | "looks" | "value"
>;

function PieRanking({
  pieRanking,
  setPieRanking,
}: {
  pieRanking: PieRankingDetails;
  setPieRanking?: (pr: PieRankingDetails) => void;
}) {
  return (
    <>
      <RankingSummary
        label="Filling"
        value={pieRanking.filling}
        setValue={
          setPieRanking &&
          ((newVal: number) =>
            setPieRanking({ ...pieRanking, filling: newVal }))
        }
      />
      <RankingSummary
        label="Pastry"
        value={pieRanking.pastry}
        setValue={
          setPieRanking &&
          ((newVal: number) => setPieRanking({ ...pieRanking, pastry: newVal }))
        }
      />
      <RankingSummary
        label="Topping"
        value={pieRanking.topping}
        setValue={
          setPieRanking &&
          ((newVal: number) =>
            setPieRanking({ ...pieRanking, topping: newVal }))
        }
      />
      <RankingSummary
        label="Looks"
        value={pieRanking.looks}
        setValue={
          setPieRanking &&
          ((newVal: number) => setPieRanking({ ...pieRanking, looks: newVal }))
        }
      />
      <RankingSummary
        label="Value"
        value={pieRanking.value}
        setValue={
          setPieRanking &&
          ((newVal: number) => setPieRanking({ ...pieRanking, value: newVal }))
        }
      />
    </>
  );
}

const validRanking = (rnk: PieRankingDetails): boolean => {
  return (
    rnk.filling > 0 &&
    rnk.pastry > 0 &&
    rnk.looks > 0 &&
    rnk.topping > 0 &&
    rnk.value > 0
  );
};

const submitRanking = (
  makerid: string,
  pieid: string,
  rankingDetails: PieRankingDetails
) => {
  return fetch(`/api/brands/${makerid}/${pieid}/ranking`, {
    method: "POST",
    body: JSON.stringify({
      ...rankingDetails,
      makerid,
      pieid,
    }),
  });
};

function SubmitPieRanking({
  makerid,
  pieid,
}: {
  makerid: string;
  pieid: string;
}) {
  const [myRanking, setMyRanking] = useState<PieRankingDetails>({
    filling: 0,
    pastry: 0,
    topping: 0,
    looks: 0,
    value: 0,
  });
  const [alreadyRanked, setAlreadyRanked] = useState<boolean>(false);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);

  const { data: session, status } = useSession();

  useEffect(() => {
    fetch(`/api/brands/${makerid}/${pieid}/ranking`, {
      method: "GET",
    })
      .then((response) => {
        if (response.status === 404) {
          setLoading(false);
        } else if (response.status === 200) {
          setAlreadyRanked(true);
          return response.json().then(setMyRanking);
        } else {
          setError(true);
        }
      })
      .catch(() => setError(true))
      .finally(() => setLoading(false));
  }, [makerid, pieid, setAlreadyRanked, setLoading, setMyRanking, setError]);

  if (status === "loading") {
    return (
      <Card>
        <CardContent>
          <Alert severity="info">Loading...</Alert>
        </CardContent>
      </Card>
    );
  }

  if (status === "unauthenticated") {
    return (
      <Card>
        <CardActions>
          <Button
            style={{ width: "100%", textAlign: "center" }}
            onClick={() => signIn()}
          >
            Sign in to Rank
          </Button>
        </CardActions>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader title="My Ranking" />
      <CardContent>
        {error && <Alert severity="error">Something has gone wrong...</Alert>}
        {loading && <Alert severity="info">Loading...</Alert>}
        {!loading && (
          <PieRanking
            pieRanking={myRanking}
            setPieRanking={!alreadyRanked ? setMyRanking : undefined}
          />
        )}
      </CardContent>
      <CardActions>
        {!alreadyRanked && (
          <Button
            style={{ width: "100%", textAlign: "center" }}
            disabled={!validRanking(myRanking) || submitting}
            onClick={() => {
              setError(false);
              setSubmitting(true);
              submitRanking(makerid, pieid, myRanking)
                .catch(() => setError(true))
                .then((res) => {
                  if (res && res.status === 200) {
                    setError(false);
                    setAlreadyRanked(true);
                  } else {
                    setError(true);
                  }
                })
                .finally(() => setSubmitting(false));
            }}
          >
            Submit
          </Button>
        )}
        <Button
          style={{ width: "100%", textAlign: "center" }}
          onClick={() => signOut()}
        >
          Log Out
        </Button>
      </CardActions>
    </Card>
  );
}

function Brands({
  maker,
  pie,
  rankingSummary,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <>
      <Head>
        <title>
          {`Mince Pie Rank :: ${maker ? maker.name : "Not Found"} :: ${
            pie ? pie.displayname : "Not Found"
          }`}
        </title>
        <meta
          name="description"
          content={`The rankings we have for ${maker?.name} ${pie?.displayname}`}
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <Breadcrumbs aria-label="breadcrumb">
          <Link color="inherit" href="/">
            Home
          </Link>
          <Link color="inherit" href={`/brands/${maker?.id}/`}>
            {maker?.name}
          </Link>
          <Typography color="text.primary">{pie?.displayname}</Typography>
        </Breadcrumbs>
        <h1>{pie && pie.displayname}</h1>
        <h3>{maker && maker.name}</h3>
        {pie && (
          <Card>
            <CardMedia
              component="img"
              height="200"
              image={pie.image_file}
              alt={`${pie.displayname}`}
            />
            <CardActions>
              <Button LinkComponent={Link} href={pie.web_link}>
                Website
              </Button>
            </CardActions>
          </Card>
        )}
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        {rankingSummary && (
          <Card>
            <CardHeader
              title="Summary"
              subheader={`${rankingSummary.count} Rankings`}
            />
            <CardContent>
              <PieRanking pieRanking={rankingSummary} />
            </CardContent>
          </Card>
        )}
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        {pie && <SubmitPieRanking makerid={pie.makerid} pieid={pie.id} />}
      </main>
    </>
  );
}

export default Brands;
