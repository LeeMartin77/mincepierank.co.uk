import { getAllRankingsForPie } from "./makerPieRankings";
import { StorageError } from "./types";

describe("Maker Storage Tests", () => {
  describe("getAllPieRankings", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getAllRankingsForPie(
        "NOT TESTED",
        "NOT TESTED",
        fakeClient as any
      );

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });
  });
});
