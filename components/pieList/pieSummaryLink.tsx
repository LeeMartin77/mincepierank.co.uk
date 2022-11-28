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
import Link from "next/link";
import { MakerPie, MakerPieRanking } from "../../system/storage";
import { ppCategory } from "../formatCategory";
import { formatPrice } from "../formatPrice";
import styles from "./pieSummaryLink.module.css";

export type PieListRanking = Omit<MakerPieRanking, "userid" | "notes"> & {
  count: number | undefined;
  average: number;
};

export function PieSummaryLink({
  pie,
  ranking,
  isTop = false,
}: {
  pie: MakerPie;
  ranking?: PieListRanking;
  isTop?: boolean;
}) {
  return (
    <Card
      data-testid={pie.id + "-link-card"}
      elevation={isTop ? 8 : 1}
      className={isTop ? styles.topCard : undefined}
    >
      <CardContent>
        <div className={styles.flexWrapper}>
          <div className={styles.flexLeftPanel}>
            <img src={pie.image_file} alt={pie.displayname} />
          </div>
          <div className={styles.flexRightPanel}>
            <Typography
              component="legend"
              className={isTop ? styles.boldText : undefined}
            >
              Avg. Ranking
            </Typography>
            <Rating
              className={isTop ? styles.topRating : undefined}
              name="read-only"
              precision={0.1}
              value={ranking?.average ?? 0}
              readOnly
            />
            {ranking && ranking.count != undefined && (
              <div className={styles.leftQuarterPad}>
                <Typography
                  variant="caption"
                  className={isTop ? styles.boldText : undefined}
                >
                  {ranking?.count ?? 0} Rankings
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
      <CardActions>
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
