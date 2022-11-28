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
      style={
        !isTop
          ? {}
          : {
              color: "black",
              fontWeight: "bold",
              background:
                "linear-gradient(29deg, rgba(255,160,27,1) 0%, rgba(255,209,144,1) 60%, rgba(252,176,69,1) 100%)",
            }
      }
    >
      <CardContent>
        <div style={{ display: "flex" }}>
          <div
            style={{
              width: "47.5%",
              height: "120px",
            }}
          >
            <img
              src={pie.image_file}
              alt={pie.displayname}
              style={{
                width: "100%",
                height: "100%",
                objectFit: "contain",
              }}
            />
          </div>
          <div
            style={{
              marginLeft: "5%",
              width: "47.5%",
            }}
          >
            <Typography
              component="legend"
              style={isTop ? { fontWeight: "bold" } : {}}
            >
              Avg. Ranking
            </Typography>
            <Rating
              style={
                isTop
                  ? {
                      boxShadow: "2px 2px 5px rgba(0,0,0,0.4)",
                      backgroundColor: "rgba(0,0,0,0.4)",
                    }
                  : {}
              }
              name="read-only"
              precision={0.1}
              value={ranking?.average ?? 0}
              readOnly
            />
            {ranking && ranking.count != undefined && (
              <div style={{ paddingLeft: "0.25em" }}>
                <Typography
                  variant="caption"
                  style={isTop ? { fontWeight: "bold" } : {}}
                >
                  {ranking?.count ?? 0} Rankings
                </Typography>
              </div>
            )}
            <div style={{ paddingLeft: "0.25em" }}>
              <Typography
                variant="caption"
                style={isTop ? { fontWeight: "bold" } : {}}
              >
                {formatPrice(pie.pack_price_in_pence / pie.pack_count)} per pie
              </Typography>
            </div>
          </div>
        </div>
        <div style={{ display: "flex", gap: "0.25em", marginTop: "0.5em" }}>
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
                  style={isTop ? { color: "black" } : {}}
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
          style={
            isTop
              ? {
                  width: "100%",
                  textAlign: "center",
                  color: "black",
                  borderColor: "black",
                  fontWeight: "bold",
                }
              : { width: "100%", textAlign: "center" }
          }
        >
          {pie.displayname}
        </Button>
      </CardActions>
    </Card>
  );
}
