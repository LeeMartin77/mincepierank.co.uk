/* eslint-disable @next/next/no-img-element */
import {
  Button,
  Card,
  CardActions,
  CardContent,
  Chip,
  Rating,
  Typography,
} from "@mui/material";
import { format } from "date-fns";
import Link from "next/link";
import { MakerPie, MakerPieRanking } from "../../system/storage";
import { ppCategory } from "../formatCategory";
import { formatPrice } from "../formatPrice";
import styles from "./pieSummaryLink.module.css";

const months = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

export type PieListRanking = Omit<MakerPieRanking, "userid" | "notes"> & {
  count: number | undefined;
  average: number;
};

export function PieSummaryLink({
  pie,
  ranking,
  showDate = false,
  isTop = false,
}: {
  pie: MakerPie;
  ranking?: PieListRanking;
  showDate?: boolean;
  isTop?: boolean;
}) {
  return (
    <Card
      data-testid={pie.id + "-link-card"}
      elevation={isTop ? 8 : 1}
      className={isTop ? styles.topCard : styles.stretch}
    >
      <CardContent>
        <div className={styles.flexWrapper}>
          <div className={styles.flexLeftPanel}>
            <img
              src={pie.image_file + "?width=300&filter=gaussian"}
              alt={pie.displayname}
            />
          </div>
          <div className={styles.flexRightPanel}>
            <Rating
              className={isTop ? styles.topRating : undefined}
              name="read-only"
              precision={0.1}
              value={ranking?.average ?? 0}
              readOnly
            />
            {ranking && ranking.last_updated && showDate && (
              <div className={styles.leftQuarterPad}>
                <Typography
                  variant="caption"
                  className={isTop ? styles.boldText : undefined}
                >
                  {format(new Date(ranking.last_updated), "MM/dd 'at' HH:mma")}
                </Typography>
              </div>
            )}
            {ranking && ranking.count != undefined && ranking.count > 0 && (
              <div className={styles.leftQuarterPad}>
                <Typography
                  variant="caption"
                  className={isTop ? styles.boldText : undefined}
                >
                  {ranking.count ?? 0} Rankings
                </Typography>
              </div>
            )}
            <div className={styles.leftQuarterPad}>
              <Typography
                variant="caption"
                className={isTop ? styles.boldText : undefined}
              >
                {formatPrice(pie.pack_price_in_pence / pie.pack_count)} per pie
              </Typography>
            </div>
          </div>
        </div>
        <div className={styles.chipContainer}>
          {pie.labels
            .sort((a, b) => a.localeCompare(b))
            .map((lb) => {
              return (
                <Chip
                  size="small"
                  key={lb}
                  label={ppCategory(lb)}
                  onClick={() => {}}
                  component={Link}
                  href={`/categories/${lb}`}
                  className={isTop ? styles.blackText : undefined}
                />
              );
            })}
        </div>
      </CardContent>
      <CardActions className={styles.actions}>
        <Button
          LinkComponent={Link}
          variant={isTop ? "outlined" : undefined}
          href={`/brands/${pie.makerid}/${pie.id}`}
          className={isTop ? styles.topButton : styles.regularButton}
        >
          {pie.displayname}
        </Button>
      </CardActions>
    </Card>
  );
}
