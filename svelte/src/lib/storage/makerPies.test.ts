import {
  getAllMakerPies,
  getPieByMakerAndId,
  getPiesByMaker,
} from "./makerPies";
import { StorageError } from "./types";

describe("Maker Storage Tests", () => {
  describe("getAllMakerPies", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getAllMakerPies(fakeClient as any);

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });
  });
  describe("getPiesByMaker", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getPiesByMaker("NOT TESTED", fakeClient as any);

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });
  });
  describe("getPieByMakerAndId", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getPieByMakerAndId(
        "NOT TESTED",
        "NOT TESTED",
        fakeClient as any
      );

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });

    test("Cassandra returns empty result set :: returns not found error", async () => {
      const fakeClient = {
        execute: jest.fn().mockResolvedValue({ rows: [] }),
      };

      const result = await getPieByMakerAndId(
        "NOT TESTED",
        "NOT TESTED",
        fakeClient as any
      );

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.NotFound);
    });
  });
});
