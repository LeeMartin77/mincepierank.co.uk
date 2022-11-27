import { Button, Divider, Grid } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import { MakerPie } from "../../system/storage";
import { PieListRanking, PieSummaryLink } from "./pieSummaryLink";

export default function InnerPieList({
  rankingOrder,
  mappedPies,
  mappedRankings,
  unrankedPies,
}: {
  mappedRankings: {
    [key: string]: PieListRanking | undefined;
  };
  mappedPies: { [key: string]: MakerPie };
  rankingOrder: string[];
  unrankedPies: Set<string>;
}) {
  const [piesToRender, setPiesToRender] = useState(9);

  useEffect(() => {
    // TODO: This is bad as it assumes the list is always the only
    // reason the page will hit the bottom, but I'm just
    // wanting to get this working
    const handleScroll = () => {
      const loadMore =
        Math.ceil(window.innerHeight + window.scrollY) >=
        document.documentElement.scrollHeight - 1500;

      if (loadMore) {
        setPiesToRender(piesToRender + 9);
      }
    };
    window.addEventListener("scroll", handleScroll, {
      passive: true,
    });

    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, [piesToRender, setPiesToRender]);

  let pieCounter = piesToRender;
  return (
    <>
      <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
      <Grid container spacing={2}>
        {rankingOrder.map((uniqid) => {
          if (pieCounter > 0) {
            pieCounter--;
            return (
              <Grid item key={uniqid} xs={12} md={6} lg={4}>
                <PieSummaryLink
                  pie={mappedPies[uniqid]}
                  ranking={mappedRankings[uniqid]}
                />
              </Grid>
            );
          }
        })}
        {Array.from(unrankedPies).map((uniqid) => {
          if (pieCounter > 0) {
            pieCounter--;
            return (
              <Grid item key={uniqid} xs={12} md={6} lg={4}>
                <PieSummaryLink pie={mappedPies[uniqid]} />
              </Grid>
            );
          }
        })}
      </Grid>
      {pieCounter === 0 && (
        <Button
          style={{ width: "100%", marginTop: "1em" }}
          variant="contained"
          onClick={() => setPiesToRender(piesToRender + 9)}
        >
          View More
        </Button>
      )}
    </>
  );
}
