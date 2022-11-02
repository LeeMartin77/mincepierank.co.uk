import { getMincePieMakers, getMincePieMaker } from "./makers";
import { StorageError } from "./types";

describe("Maker Storage Tests", () => {
  describe("getMincePieMakers", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getMincePieMakers(fakeClient as any);

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });
  });

  describe("getMincePieMaker", () => {
    test("Cassandra throws error :: returns general error", async () => {
      const fakeClient = {
        execute: jest.fn().mockRejectedValue({}),
      };

      const result = await getMincePieMaker("NOT TESTED", fakeClient as any);

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.GenericError);
    });

    test("Cassandra returns empty result set :: returns not found error", async () => {
      const fakeClient = {
        execute: jest.fn().mockResolvedValue({ rows: [] }),
      };

      const result = await getMincePieMaker("NOT TESTED", fakeClient as any);

      expect(result.isErr()).toBeTruthy();
      expect(result._unsafeUnwrapErr()).toBe(StorageError.NotFound);
    });
  });
});
