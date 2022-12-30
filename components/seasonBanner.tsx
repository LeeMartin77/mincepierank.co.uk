import { getNextSeasonEndDate, timeInSeason } from "../system/seasonConfig";
import styles from "./seasonBanner.module.css";

// Seven days in ms
const WARNING_TIME_LENGTH = 1000 * 60 * 60 * 24 * 7;

export function SeasonBanner({
  nowFn = Date.now,
  nextEndDateFn = getNextSeasonEndDate,
}) {
  const now = nowFn();
  if (!timeInSeason(now)) {
    return (
      <div
        className={styles.banner}
        style={{ backgroundColor: "#3c66c1" }}
        data-testid="out-of-season-banner"
      >
        Mince Pie Rank is currently closed until the next Christmas period!
      </div>
    );
  }
  if (nextEndDateFn() - WARNING_TIME_LENGTH < now) {
    return (
      <div
        className={styles.banner}
        style={{ backgroundColor: "#cc5730" }}
        data-testid="near-season-end-banner"
      >
        Mince Pie Rank season is nearly over! Ending{" "}
        {new Date(nextEndDateFn()).toLocaleString()}
      </div>
    );
  }

  return <></>;
}
