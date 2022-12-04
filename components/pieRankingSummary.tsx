import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Grid,
  Rating,
  Typography,
} from "@mui/material";
import { ExpandMore } from "@mui/icons-material";
import { PieRankingSummary } from "../system/storage";
import { calculateAverage } from "../system/storage/utilities";
import { useState } from "react";
import styles from "./pieRankingSummary.module.css";

function RankingSummary({ label, value }: { label: string; value: number }) {
  const xs = 12;
  const sm = 6;
  const md = 4;
  return (
    <div className={styles.rankingItem}>
      <Rating name="read-only" value={value} readOnly={true} precision={0.1} />
      <Typography component="legend" textAlign={"center"}>
        {`${label} (${value.toFixed(1)})`}
      </Typography>
    </div>
  );
}

function PieRanking({ pieRanking }: { pieRanking: PieRankingSummary }) {
  return (
    <div className={styles.rankingContainer}>
      <RankingSummary label="Filling" value={pieRanking.filling} />
      <RankingSummary label="Pastry" value={pieRanking.pastry} />
      <RankingSummary label="Topping" value={pieRanking.topping} />
      <RankingSummary label="Looks" value={pieRanking.looks} />
      <RankingSummary label="Value" value={pieRanking.value} />
    </div>
  );
}

export function PieRankingSummary({
  pieRanking,
}: {
  pieRanking: PieRankingSummary;
}) {
  const [expanded, setExpanded] = useState(false);
  return (
    <Accordion
      expanded={expanded}
      onChange={() => setExpanded(!expanded)}
      elevation={1}
    >
      <AccordionSummary
        expandIcon={<ExpandMore />}
        aria-controls="summary-content"
      >
        <div
          style={{
            display: "flex",
            flexDirection: "row",
            justifyContent: "space-around",
            width: "100%",
          }}
        >
          <div className={styles.rankingItem}>
            <Rating
              size="large"
              name="read-only"
              value={calculateAverage(pieRanking)}
              readOnly={true}
              precision={0.1}
            />
            <Typography
              component="legend"
              textAlign={"center"}
              fontWeight={"bold"}
            >
              {`${calculateAverage(pieRanking).toFixed(1)}/5 with ${
                pieRanking.count
              } votes`}
            </Typography>
          </div>
        </div>
      </AccordionSummary>
      <AccordionDetails>
        <PieRanking pieRanking={pieRanking} />
      </AccordionDetails>
    </Accordion>
  );
}
