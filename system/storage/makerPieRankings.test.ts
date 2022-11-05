import {
  getAllRankingsForPie,
  getPieRankingSummary,
  addPieRanking,
} from "./makerPieRankings";
import { MakerPieRanking, StorageError } from "./types";
import { randomBytes } from "node:crypto";

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

  describe("addPieRanking", () => {
    const badRankings = [
      {
        makerid: "test-maker-id",
        pieid: "test-pie-id",
        userid: "test-user-id",
        pastry: 0,
        filling: 0,
        topping: 0,
        looks: 0,
        value: 0,
      },
      {
        makerid: "test-maker-id",
        pieid: "test-pie-id",
        userid: "test-user-id",
        pastry: 6,
        filling: 7,
        topping: 1,
        looks: 2,
        value: 3,
      },
      {
        pastry: 4,
        filling: 2,
        topping: 1,
        looks: 2,
        value: 3,
      } as MakerPieRanking,
      {
        makerid: "test-maker-id",
        pieid: "test-pie-id",
        userid: "test-user-id",
        pastry: 1,
        filling: 2,
        topping: 3,
        looks: 4,
        value: 5,
        notes: "Longer than 140 Char" + randomBytes(140).toString("hex"),
      },
    ];
    test.each(badRankings)(
      "Invalid input :: Rejects with error",
      async (ranking: MakerPieRanking) => {
        const fakeClient = {
          execute: jest.fn().mockRejectedValue({}),
        };

        const result = await addPieRanking(ranking);

        expect(fakeClient.execute).not.toBeCalled();
        expect(result.isErr()).toBeTruthy();
        expect(result._unsafeUnwrapErr()).toBe(StorageError.BadInput);
      }
    );
    test("Valid Input :: Client Rejects :: Handles Gracefully", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const ranking: MakerPieRanking = {
        makerid: "test-maker-id",
        pieid: "test-pie-id",
        userid: "test-user-id",
        pastry: 4,
        filling: 5,
        topping: 1,
        looks: 2,
        value: 3,
      };

      const result = await addPieRanking(ranking, fakeClient as any);

      expect(fakeClient.execute).toBeCalled();
      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });
  });
});
