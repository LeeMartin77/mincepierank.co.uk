import { Grid, Rating, Typography } from "@mui/material";
import { PieRankingSummary } from "../system/storage";
import { calculateAverage } from "../system/storage/utilities";

function RankingSummary({
  label,
  value,
  bold = false,
}: {
  label: string;
  value: number;
  bold?: boolean;
}) {
  return (
    <>
      <Typography component="legend" style={bold ? { fontWeight: "bold" } : {}}>
        {`${label} (${value.toFixed(1)})`}
      </Typography>
      <Rating name="read-only" value={value} readOnly={true} precision={0.1} />
    </>
  );
}

function PieRanking({ pieRanking }: { pieRanking: PieRankingSummary }) {
  const xs = 12;
  const sm = 6;
  const md = 4;
  return (
    <Grid container spacing={2}>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankingSummary label="Filling" value={pieRanking.filling} />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankingSummary label="Pastry" value={pieRanking.pastry} />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankingSummary label="Topping" value={pieRanking.topping} />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankingSummary label="Looks" value={pieRanking.looks} />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankingSummary label="Value" value={pieRanking.value} />
      </Grid>
      <Grid item xs={xs} sm={sm} md={md}>
        <RankingSummary
          label="Total Score"
          value={calculateAverage(pieRanking)}
          bold
        />
      </Grid>
    </Grid>
  );
}

export function PieRankingSummary({
  pieRanking,
}: {
  pieRanking: PieRankingSummary;
}) {
  return <PieRanking pieRanking={pieRanking} />;
}
