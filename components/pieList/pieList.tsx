/* eslint-disable @next/next/no-img-element */
import { MakerPie } from "../../system/storage";
import { Chip, Divider, Typography } from "@mui/material";
import { Check } from "@mui/icons-material";
import { mapPiesAndRankings } from "../mapPiesAndRankings";
import Head from "next/head";
import { descriptionSummary } from "../descriptionSummary";
import { useState } from "react";
import { ppCategory } from "../formatCategory";
import { PieListRanking, PieSummaryLink } from "./pieSummaryLink";
import InnerPieList from "./innerPieList";

function PieFilter({
  availableCategories,
  filteredCategories,
  setFilteredCategories,
}: {
  availableCategories: Set<string>;
  filteredCategories: Set<string>;
  setFilteredCategories: (cats: Set<string>) => void;
}) {
  return (
    <div
      style={{
        display: "flex",
        gap: "0.5em",
        maxWidth: "100%",
        scrollbarWidth: "none",
        overflowX: "auto",
      }}
    >
      {Array.from(availableCategories)
        .sort((a, b) => a.localeCompare(b))
        .map((available) => {
          if (filteredCategories.has(available)) {
            return (
              <Chip
                key={available}
                icon={<Check />}
                label={ppCategory(available)}
                onClick={() => {
                  const newSet = new Set(filteredCategories);
                  newSet.delete(available);
                  setFilteredCategories(newSet);
                }}
              />
            );
          }
          return (
            <Chip
              key={available}
              label={ppCategory(available)}
              variant="outlined"
              onClick={() => {
                const newSet = new Set(filteredCategories);
                newSet.add(available);
                setFilteredCategories(newSet);
              }}
            />
          );
        })}
    </div>
  );
}

export function PieList({
  pies,
  rankings,
  addMetaDescription = false,
  metaPrefix = "",
  lockedCategory = undefined,
  showDates = false,
}: {
  pies: MakerPie[];
  rankings: PieListRanking[];
  addMetaDescription?: boolean;
  metaPrefix?: string;
  lockedCategory?: string;
  showDates?: boolean;
}) {
  const availableCategories = pies.reduce((prev, pie) => {
    return new Set([...Array.from(prev), ...pie.labels]);
  }, new Set([] as string[]));
  if (lockedCategory) {
    availableCategories.delete(lockedCategory);
  }
  const [filteredCategories, setFilteredCategories] = useState(
    new Set([]) as Set<string>
  );

  const { mappedRankings, mappedPies, rankingOrder, unrankedPies } =
    mapPiesAndRankings(pies, rankings, filteredCategories);
  const topPieId = rankingOrder.shift();
  const showFurtherList = rankingOrder.length > 0 || unrankedPies.size > 0;
  return (
    <>
      {addMetaDescription && (
        <Head>
          <meta
            name="description"
            content={`${metaPrefix}${
              topPieId
                ? descriptionSummary(
                    mappedPies[topPieId]!,
                    mappedRankings[topPieId]!
                  )
                : ""
            }`}
          />
        </Head>
      )}
      <PieFilter
        availableCategories={availableCategories}
        filteredCategories={filteredCategories}
        setFilteredCategories={setFilteredCategories}
      />
      <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />

      {topPieId && mappedPies[topPieId] && (
        <>
          <h1>Top Pie</h1>
          <PieSummaryLink
            isTop
            pie={mappedPies[topPieId]}
            ranking={mappedRankings[topPieId]}
            showDate={showDates}
          />
        </>
      )}
      {!topPieId && unrankedPies.size > 0 && (
        <Typography
          variant="h4"
          style={{
            width: "100%",
            marginTop: "1em",
            marginBottom: "1em",
            textAlign: "center",
          }}
        >
          No top pie yet...
        </Typography>
      )}
      {showFurtherList && (
        <InnerPieList
          mappedPies={mappedPies}
          mappedRankings={mappedRankings}
          rankingOrder={rankingOrder}
          unrankedPies={unrankedPies}
          showDates={showDates}
        />
      )}
      {!showFurtherList && !topPieId && (
        <Typography
          variant="h4"
          style={{
            width: "100%",
            marginTop: "1em",
            marginBottom: "1em",
            textAlign: "center",
          }}
        >
          No pies meeting criteria - try removing some categories.
        </Typography>
      )}
    </>
  );
}
