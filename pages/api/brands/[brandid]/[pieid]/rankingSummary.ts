// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from "next";
import {
  StorageError,
  getPieRankingSummary,
} from "../../../../../system/storage";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (!["GET"].includes(req.method as string)) {
    return res.status(400).send("GET Requests only");
  }
  const { brandid, pieid } = req.query;
  const summary = await getPieRankingSummary(
    brandid as string,
    pieid as string
  );
  summary
    .map((summary) => {
      res.status(200).send(summary);
    })
    .mapErr((err) => {
      res.status(StorageError.NotFound === err ? 404 : 500).send("Error");
    });
}
