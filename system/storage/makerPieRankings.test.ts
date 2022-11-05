import { getAllRankingsForPie, getPieRankingSummary } from "./makerPieRankings";
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

  describe("getPieRankingSummary", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getPieRankingSummary(
        "NOT TESTED",
        "NOT TESTED",
        fakeClient as any
      );

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });

    test("Cassandra returns results :: aggregates results", async () => {
      const makerid = "test-maker-id";
      const pieid = "test-pie-id";
      const fakeClient = {
        execute: jest.fn().mockResolvedValue({
          rows: [
            {
              makerid,
              pieid,
              pastry: 2,
              filling: 2,
              topping: 2,
              looks: 2,
              value: 2,
            },
            {
              makerid,
              pieid,
              pastry: 3,
              filling: 3,
              topping: 3,
              looks: 3,
              value: 3,
            },
            {
              makerid,
              pieid,
              pastry: 4,
              filling: 4,
              topping: 4,
              looks: 4,
              value: 4,
            },
          ],
        }),
      };

      const fakeMapper = (val: any) => val;

      const result = await getPieRankingSummary(
        makerid,
        pieid,
        fakeClient as any,
        fakeMapper
      );

      expect(fakeClient.execute).toBeCalledTimes(1);
      const call = fakeClient.execute.mock.calls[0];

      expect(call[1]).toEqual([makerid, pieid]);

      expect(result.isOk()).toBeTruthy();
      expect(result._unsafeUnwrap()).toEqual({
        count: 3,
        filling: 3,
        looks: 3,
        makerid: "test-maker-id",
        pastry: 3,
        pieid: "test-pie-id",
        topping: 3,
        value: 3,
      });
    });
  });
});
