import { err, ok } from "neverthrow";
import { StorageError } from "../../../../../system/storage";
import handler from "./ranking";

const testEmail = "happypath@email.com";

jest.mock("next-auth/jwt", () => {
  return {
    getToken: jest.fn().mockResolvedValue({ id: "happypath@email.com" }),
  };
});

const happypathDatefunction = () => new Date(2022, 11, 12).getTime();

describe("Pie Ranking Endpoint", () => {
  const badMethods = ["PUT", "DELETE"];

  test("Not authenticated :: reject with 401", async () => {
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      send: jest.fn(),
    };
    await handler(
      {
        method: "GET",
        body: JSON.stringify({ makerid: "something", pieid: "somethingelse" }),
      } as any,
      mockResponse as any,
      undefined,
      undefined,
      undefined,
      jest.fn().mockResolvedValue(null)
    );
    expect(mockResponse.status).toBeCalledWith(401);
    expect(mockResponse.send).toBeCalled();
    expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
  });

  test.each(badMethods)("%s Request :: Returns Bad Result", (method) => {
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      send: jest.fn(),
    };
    handler(
      {
        method,
        body: JSON.stringify({ makerid: "something", pieid: "somethingelse" }),
      } as any,
      mockResponse as any
    );
    expect(mockResponse.status).toBeCalledWith(400);
    expect(mockResponse.send).toBeCalled();
    expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
  });
  describe("GET request", () => {
    test("Happy Path :: Passes Parameters and returns result", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const mockQuery = { brandid: "test-brand-id", pieid: "test-pie-id" };
      const mockCassResponse = { something: "whatever" };
      const fakeMethod = jest.fn().mockResolvedValue(ok(mockCassResponse));
      await handler(
        {
          method: "GET",
          query: mockQuery,
        } as any,
        mockResponse as any,
        jest.fn(),
        jest.fn(),
        fakeMethod
      );
      expect(fakeMethod).toBeCalledWith(
        mockQuery.brandid,
        mockQuery.pieid,
        testEmail
      );
      expect(mockResponse.status).toBeCalledWith(200);
      expect(mockResponse.send).toBeCalledWith(mockCassResponse);
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });
    test("Storage Error :: Returns 500", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const mockQuery = { brandid: "test-brand-id", pieid: "test-pie-id" };
      const fakeMethod = jest
        .fn()
        .mockResolvedValue(err(StorageError.GenericError));
      await handler(
        {
          method: "GET",
          query: mockQuery,
        } as any,
        mockResponse as any,
        jest.fn(),
        jest.fn(),
        fakeMethod
      );
      expect(mockResponse.status).toBeCalledWith(500);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });
    test("Not Found :: Returns 404", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const mockQuery = { brandid: "test-brand-id", pieid: "test-pie-id" };
      const fakeMethod = jest
        .fn()
        .mockResolvedValue(err(StorageError.NotFound));
      await handler(
        {
          method: "GET",
          query: mockQuery,
        } as any,
        mockResponse as any,
        jest.fn(),
        jest.fn(),
        fakeMethod
      );
      expect(mockResponse.status).toBeCalledWith(404);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });
  });
  describe("POST request", () => {
    test("Missing makerid and pieid :: Returns 400", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const getPieFunction = jest.fn();
      const insertStorageFunction = jest.fn();
      const fakeBody = { something: "else" };
      await handler(
        { method: "POST", body: JSON.stringify(fakeBody) } as any,
        mockResponse as any,
        getPieFunction,
        insertStorageFunction,
        undefined,
        undefined,
        happypathDatefunction
      );
      expect(getPieFunction).not.toBeCalled();
      expect(insertStorageFunction).not.toBeCalled();
      expect(mockResponse.status).toBeCalledWith(400);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });

    test("Pie Not Found :: Returns 400", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const makerid = "test-maker-id";
      const pieid = "test-pie-id";
      const getPieFunction = jest
        .fn()
        .mockResolvedValue(err(StorageError.NotFound));
      const insertStorageFunction = jest.fn();
      const fakeBody = { makerid, pieid, something: "else" };
      await handler(
        { method: "POST", body: JSON.stringify(fakeBody) } as any,
        mockResponse as any,
        getPieFunction,
        insertStorageFunction,
        undefined,
        undefined,
        happypathDatefunction
      );
      expect(getPieFunction).toBeCalledWith(makerid, pieid);
      expect(insertStorageFunction).not.toBeCalled();
      expect(mockResponse.status).toBeCalledWith(400);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });

    test("Pie Storage Error :: Returns 500", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const makerid = "test-maker-id";
      const pieid = "test-pie-id";
      const getPieFunction = jest
        .fn()
        .mockResolvedValue(err(StorageError.GenericError));
      const insertStorageFunction = jest.fn();
      const fakeBody = { makerid, pieid, something: "else" };
      await handler(
        { method: "POST", body: JSON.stringify(fakeBody) } as any,
        mockResponse as any,
        getPieFunction,
        insertStorageFunction,
        undefined,
        undefined,
        happypathDatefunction
      );
      expect(getPieFunction).toBeCalledWith(makerid, pieid);
      expect(insertStorageFunction).not.toBeCalled();
      expect(mockResponse.status).toBeCalledWith(500);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });

    test("Passes Request :: Bad Request Response :: Returns 400", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const getPieFunction = jest.fn().mockResolvedValue(ok({}));
      const insertStorageFunction = jest
        .fn()
        .mockResolvedValue(err(StorageError.BadInput));
      const fakeBody = {
        makerid: "test-maker-id",
        pieid: "test-pie-id",
        something: "else",
      };
      await handler(
        { method: "POST", body: JSON.stringify(fakeBody) } as any,
        mockResponse as any,
        getPieFunction,
        insertStorageFunction,
        undefined,
        undefined,
        happypathDatefunction
      );
      expect(insertStorageFunction).toBeCalledWith({
        ...fakeBody,
        userid: testEmail,
      });
      expect(mockResponse.status).toBeCalledWith(400);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });

    test("Passes Request :: Storage Error Response :: Returns 500", async () => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const getPieFunction = jest.fn().mockResolvedValue(ok({}));
      const insertStorageFunction = jest
        .fn()
        .mockResolvedValue(err(StorageError.GenericError));
      const fakeBody = {
        makerid: "test-maker-id",
        pieid: "test-pie-id",
        something: "else",
      };
      await handler(
        { method: "POST", body: JSON.stringify(fakeBody) } as any,
        mockResponse as any,
        getPieFunction,
        insertStorageFunction,
        undefined,
        undefined,
        happypathDatefunction
      );
      expect(insertStorageFunction).toBeCalledWith({
        ...fakeBody,
        userid: testEmail,
      });
      expect(mockResponse.status).toBeCalledWith(500);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    });

    const goodDates = [
      new Date(2022, 11, 31, 23, 58, 58),
      new Date(2022, 10, 12, 12, 30, 30),
      new Date(2022, 9, 1, 0, 1, 1),
    ];

    test.each(goodDates)(
      "Passes Request :: Good Response :: Returns 200 OK",
      async (date) => {
        const mockResponse = {
          status: jest.fn().mockReturnThis(),
          send: jest.fn(),
        };
        const getPieFunction = jest.fn().mockResolvedValue(ok({}));
        const insertStorageFunction = jest.fn().mockResolvedValue(ok(true));
        const fakeBody = {
          makerid: "test-maker-id",
          pieid: "test-pie-id",
          something: "else",
        };
        await handler(
          { method: "POST", body: JSON.stringify(fakeBody) } as any,
          mockResponse as any,
          getPieFunction,
          insertStorageFunction,
          undefined,
          undefined,
          () => date.getTime()
        );
        expect(insertStorageFunction).toBeCalledWith({
          ...fakeBody,
          userid: testEmail,
        });
        expect(mockResponse.status).toBeCalledWith(200);
        expect(mockResponse.send).toBeCalledWith({ status: "ok" });
        expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
      }
    );
  });

  const badDates = [
    new Date(2023, 0, 1, 0, 1, 0),
    new Date(2022, 8, 30, 22, 58, 0),
    new Date(2022, 4, 12, 12, 30, 0),
  ];

  test.each(badDates)(
    "POST Outside of season :: %s :: rejects with 400",
    async (date) => {
      const mockResponse = {
        status: jest.fn().mockReturnThis(),
        send: jest.fn(),
      };
      const getPieFunction = jest.fn();
      const insertStorageFunction = jest.fn();
      const fakeBody = {
        makerid: "test-maker-id",
        pieid: "test-pie-id",
        something: "else",
      };
      await handler(
        { method: "POST", body: JSON.stringify(fakeBody) } as any,
        mockResponse as any,
        getPieFunction,
        insertStorageFunction,
        undefined,
        undefined,
        () => date.getTime()
      );
      expect(getPieFunction).not.toBeCalled();
      expect(insertStorageFunction).not.toBeCalled();
      expect(mockResponse.status).toBeCalledWith(400);
      expect(mockResponse.send).toBeCalled();
      expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
    }
  );
});
