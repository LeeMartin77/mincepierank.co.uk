import {
  getAllRankingsForPie,
  getPieRankingSummary,
  addPieRanking,
  getMyRankingForPie,
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

  describe("getMyRankingForPie", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getMyRankingForPie(
        "NOT TESTED",
        "NOT TESTED",
        "NOT TESTED",
        fakeClient as any
      );

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });
    test("Found nothing in DB :: returns not found error", async () => {
      const fakeClient = {
        execute: jest.fn().mockResolvedValue({
          rows: [],
        }),
      };

      const result = await getMyRankingForPie(
        "NOT TESTED",
        "NOT TESTED",
        "NOT TESTED",
        fakeClient as any
      );

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.NotFound);
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
