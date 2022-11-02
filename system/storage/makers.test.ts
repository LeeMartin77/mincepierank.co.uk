import { getMincePieMakers } from "./makers";
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
});
