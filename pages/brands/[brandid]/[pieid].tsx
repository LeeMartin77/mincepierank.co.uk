import Head from "next/head";
import { InferGetServerSidePropsType } from "next";
import {
  getMincePieMaker,
  getPieByMakerAndId,
  getPieRankingSummary,
  MakerPie,
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
  Chip,
  Divider,
  Grid,
  IconButton,
  Rating,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import { Link as LinkIcon } from "@mui/icons-material";
import { useCallback, useEffect, useState } from "react";
import { signIn, useSession } from "next-auth/react";
import Link from "next/link";
import { formatPrice } from "../../../components/formatPrice";
import { Share } from "@mui/icons-material";
import { calculateAverage } from "../../../system/storage/utilities";
import { ppCategory } from "../../../components/formatCategory";
import { PieRankingSummary } from "../../../components/pieRankingSummary";
import { format } from "date-fns";

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

function RankPieAttribute({
  label,
  value,
  setValue,
}: {
  label: string;
  value: number;
  setValue: (n: number) => void;
}) {
  return (
    <>
      <Typography component="legend">
        {`${label}${!setValue ? ` (${value.toFixed(1)})` : ``}`}
      </Typography>
      <Rating
        name={label}
        value={value}
        precision={1}
        onChange={(event, newValue) => {
          newValue && setValue(newValue);
        }}
      />
    </>
  );
}

type PieRankingDetails = Pick<
  MakerPieRanking,
  "filling" | "pastry" | "topping" | "looks" | "value" | "last_updated"
>;

function UserRankingControl({
  pieRanking,
  setPieRanking,
}: {
  pieRanking: PieRankingDetails;
  setPieRanking: (pr: PieRankingDetails) => void;
}) {
  const xs = 12;
  const sm = 6;
  const md = 4;
  return (
    <Grid container spacing={2}>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankPieAttribute
          label="Filling"
          value={pieRanking.filling}
          setValue={(newVal: number) =>
            setPieRanking({ ...pieRanking, filling: newVal })
          }
        />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankPieAttribute
          label="Pastry"
          value={pieRanking.pastry}
          setValue={(newVal: number) =>
            setPieRanking({ ...pieRanking, pastry: newVal })
          }
        />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankPieAttribute
          label="Topping"
          value={pieRanking.topping}
          setValue={(newVal: number) =>
            setPieRanking({ ...pieRanking, topping: newVal })
          }
        />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankPieAttribute
          label="Looks"
          value={pieRanking.looks}
          setValue={(newVal: number) =>
            setPieRanking({ ...pieRanking, looks: newVal })
          }
        />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankPieAttribute
          label="Value"
          value={pieRanking.value}
          setValue={(newVal: number) =>
            setPieRanking({ ...pieRanking, value: newVal })
          }
        />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <Typography component="legend" style={{ fontWeight: "bold" }}>
          {`Total Score  (${calculateAverage(pieRanking).toFixed(1)})`}
        </Typography>
        <Rating
          name={"summary"}
          value={calculateAverage(pieRanking)}
          precision={0.1}
          readOnly
        />
      </Grid>
    </Grid>
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

function rankingsDifferent(
  rankingOne: { [key: string]: any },
  rankingTwo: { [key: string]: any }
) {
  const fields = ["filling", "pastry", "topping", "looks", "value"];
  for (let field of fields) {
    if (rankingOne[field] != rankingTwo[field]) {
      return true;
    }
  }
  return false;
}

function SubmitPieRanking({
  pie,
  setRefresh,
}: {
  pie: MakerPie;
  setRefresh: (inp: boolean) => void;
}) {
  const [myRanking, setMyRanking] = useState<PieRankingDetails>({
    filling: 0,
    pastry: 0,
    topping: 0,
    looks: 0,
    value: 0,
  });

  const { makerid, id: pieid } = pie;

  const [alreadyRanked, setAlreadyRanked] = useState<
    PieRankingDetails | undefined
  >(undefined);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);

  const saveRanking = useCallback(
    (imakerid: string, ipieid: string, iranking: PieRankingDetails) => {
      setError(false);
      setSubmitting(true);
      submitRanking(imakerid, ipieid, iranking)
        .catch(() => setError(true))
        .then((res) => {
          if (res && res.status === 200) {
            setError(false);
            setAlreadyRanked(iranking);
            setRefresh(true);
          } else {
            setError(true);
          }
        })
        .finally(() => setSubmitting(false));
    },
    [setError, setSubmitting, setAlreadyRanked, setRefresh]
  );

  useEffect(() => {
    if (
      validRanking(myRanking) &&
      (!alreadyRanked || rankingsDifferent(alreadyRanked, myRanking))
    ) {
      saveRanking(makerid, pieid, myRanking);
    }
  }, [saveRanking, alreadyRanked, myRanking, makerid, pieid]);

  const { data: session, status } = useSession();

  useEffect(() => {
    fetch(`/api/brands/${makerid}/${pieid}/ranking`, {
      method: "GET",
    })
      .then((response) => {
        if (response.status === 404) {
          setLoading(false);
        } else if (response.status === 200) {
          return response.json().then((res) => {
            res.last_updated = res.last_updated
              ? new Date(res.last_updated)
              : undefined;
            setAlreadyRanked(res);
            setMyRanking(res);
          });
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
        <CardHeader title="Sign in to Rank" />
        <CardActions>
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
        </CardActions>
      </Card>
    );
  }

  let subHeader = "Loading...";

  if (!submitting && alreadyRanked) {
    if (myRanking.last_updated) {
      subHeader =
        "Ranked on " +
        format(new Date(myRanking.last_updated), "yyyy/MM/dd 'at' HH:mm");
    } else {
      subHeader = "Ranked!";
    }
  }

  if (!submitting && !alreadyRanked) {
    subHeader = "Awaiting Ranking";
  }

  if (submitting && !alreadyRanked) {
    subHeader = "Submitting...";
  }

  if (submitting && alreadyRanked) {
    subHeader = "Updating...";
  }

  return (
    <Card>
      <CardHeader
        title="My Ranking"
        subheader={subHeader}
        action={
          alreadyRanked && navigator.share ? (
            <IconButton
              aria-label="share"
              onClick={() =>
                navigator.share({
                  url: `https://mincepierank.co.uk/brands/${pie.makerid}/${pie.id}`,
                  title: `${pie.displayname} :: Mince Pie Rank`,
                  text: `${pie.displayname} :: I gave it ${calculateAverage(
                    myRanking
                  ).toFixed(1)}/5 on Mince Pie Rank`,
                })
              }
            >
              <Share />
            </IconButton>
          ) : undefined
        }
      />
      <CardContent>
        {error && <Alert severity="error">Something has gone wrong...</Alert>}
        {loading && <Alert severity="info">Loading...</Alert>}
        {!loading && (
          <UserRankingControl
            pieRanking={myRanking}
            setPieRanking={setMyRanking}
          />
        )}
      </CardContent>
    </Card>
  );
}

function PieDetails({ pie }: { pie: MakerPie }) {
  const statistics: { name: string; value: string }[] = [
    { name: "Pack Quantity", value: pie.pack_count.toString() },
    { name: "Pack Price", value: formatPrice(pie.pack_price_in_pence) },
    {
      name: "Price per Pie",
      value: formatPrice(pie.pack_price_in_pence / pie.pack_count),
    },
  ];
  return (
    <Card>
      <CardMedia
        component="img"
        height="200"
        image={pie.image_file}
        alt={`${pie.displayname}`}
        style={{
          objectFit: "contain",
        }}
      />
      <div style={{ display: "flex", gap: "0.5em", margin: "0.5em 1em" }}>
        {pie.labels
          .sort((a, b) => a.localeCompare(b))
          .map((lb) => {
            return (
              <Chip
                key={lb}
                label={ppCategory(lb)}
                onClick={() => {}}
                component={Link}
                href={`/categories/${lb}`}
              />
            );
          })}
      </div>
      <TableContainer component={CardContent}>
        <Table size="small" aria-label="Mince Pie Stats">
          <TableHead>
            <TableRow>
              <TableCell>Statistics</TableCell>
              <TableCell></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {statistics.map((val, i) => (
              <TableRow
                key={`property-${i}`}
                sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
              >
                <TableCell component="th" scope="row">
                  {val.name}
                </TableCell>
                <TableCell>{val.value}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <CardActions>
        <Button
          variant="outlined"
          LinkComponent={Link}
          href={pie.web_link}
          style={{ width: "100%" }}
          endIcon={<LinkIcon />}
        >
          Go To Website
        </Button>
      </CardActions>
    </Card>
  );
}

function Pie({
  maker,
  pie,
  rankingSummary,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const [refresh, setRefresh] = useState(false);
  const [localSummary, setLocalSummary] = useState(rankingSummary);
  useEffect(() => {
    if (refresh) {
      fetch(`/api/brands/${maker.id}/${pie.id}/rankingSummary`).then((res) => {
        res.json().then((summary) => {
          setLocalSummary(summary);
          setRefresh(false);
        });
      });
    }
  }, [maker.id, pie.id, refresh, setLocalSummary, setRefresh]);

  useEffect(() => {
    if (navigator.canShare) {
      navigator.canShare({
        url: `https://mincepierank.co.uk/brands/${pie.makerid}/${pie.id}`,
        title: `${pie.displayname} :: Mince Pie Rank`,
        text: `${pie.displayname} has ${
          rankingSummary?.average.toFixed(1) ?? 0
        }/5 on Mince Pie Rank`,
      });
    }
  });
  return (
    <>
      <Head>
        <title>{`${
          rankingSummary && calculateAverage(rankingSummary) > 0
            ? `${calculateAverage(rankingSummary)}/5`
            : "No rankings"
        } for ${pie.displayname} from ${maker.name} :: Mince Pie Rank`}</title>
        <meta
          name="description"
          content={`${pie.displayname} from ${maker.name} has ${
            rankingSummary?.average.toFixed(1) ?? 0
          } stars across ${rankingSummary?.count ?? 0} votes.`}
        />
        <meta property="og:image" content="/logo-social.png" />
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
        <h1>{pie.displayname}</h1>
        <h3>{maker.name}</h3>
        <PieDetails pie={pie} />
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        {localSummary && <PieRankingSummary pieRanking={localSummary} />}
        <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
        {pie && <SubmitPieRanking pie={pie} setRefresh={setRefresh} />}
      </main>
    </>
  );
}

export default Pie;
