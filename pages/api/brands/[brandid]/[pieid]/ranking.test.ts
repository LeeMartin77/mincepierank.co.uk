import { err, ok } from "neverthrow";
import { StorageError } from "../../../../../system/storage";
import handler from "./ranking";

describe("Pie Ranking Endpoint", () => {
  const badMethods = ["GET", "PUT", "DELETE"];

  test.each(badMethods)("%s Request :: Returns Bad Result", (method) => {
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      send: jest.fn(),
    };
    handler({ method } as any, mockResponse as any);
    expect(mockResponse.status).toBeCalledWith(400);
    expect(mockResponse.send).toBeCalled();
    expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
  });

  test("Missing makerid and pieid :: Returns 400", async () => {
    const mockResponse = {
      status: jest.fn().mockReturnThis(),
      send: jest.fn(),
    };
    const getPieFunction = jest.fn();
    const insertStorageFunction = jest.fn();
    const fakeBody = { something: "else" };
    await handler(
      { method: "POST", body: fakeBody } as any,
      mockResponse as any,
      getPieFunction,
      insertStorageFunction
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
      { method: "POST", body: fakeBody } as any,
      mockResponse as any,
      getPieFunction,
      insertStorageFunction
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
      { method: "POST", body: fakeBody } as any,
      mockResponse as any,
      getPieFunction,
      insertStorageFunction
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
      { method: "POST", body: fakeBody } as any,
      mockResponse as any,
      getPieFunction,
      insertStorageFunction
    );
    expect(insertStorageFunction).toBeCalledWith(fakeBody);
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
      { method: "POST", body: fakeBody } as any,
      mockResponse as any,
      getPieFunction,
      insertStorageFunction
    );
    expect(insertStorageFunction).toBeCalledWith(fakeBody);
    expect(mockResponse.status).toBeCalledWith(500);
    expect(mockResponse.send).toBeCalled();
    expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
  });

  test("Passes Request :: Good Response :: Returns 200 OK", async () => {
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
      { method: "POST", body: fakeBody } as any,
      mockResponse as any,
      getPieFunction,
      insertStorageFunction
    );
    expect(insertStorageFunction).toBeCalledWith(fakeBody);
    expect(mockResponse.status).toBeCalledWith(200);
    expect(mockResponse.send).toBeCalledWith({ status: "ok" });
    expect(mockResponse.status).toHaveBeenCalledBefore(mockResponse.send);
  });
});
