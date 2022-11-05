// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from "next";
import {
  addPieRanking,
  getPieByMakerAndId,
  StorageError,
} from "../../../../../system/storage";

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<{ status: "ok" | "ko" } | string>,
  fnGetPie = getPieByMakerAndId,
  fnAddPieRanking = addPieRanking
) {
  //const { brandid, pieid } = req.query;
  const parsedBody = JSON.parse(req.body);
  if (req.method === "POST" && parsedBody.makerid && parsedBody.pieid) {
    return fnGetPie(parsedBody.makerid, parsedBody.pieid).then((pieRes) => {
      pieRes
        .mapErr((se) => {
          if (se === StorageError.NotFound) {
            res.status(400).send("Bad Request");
          } else {
            res.status(500).send("Server Error");
          }
        })
        .map(() => {
          fnAddPieRanking(parsedBody).then((sres) => {
            sres
              .map(() => {
                res.status(200).send({ status: "ok" });
              })
              .mapErr((se) => {
                if (se === StorageError.BadInput) {
                  res.status(400).send("Bad Request");
                } else {
                  res.status(500).send("Server Error");
                }
              });
          });
        });
    });
  } else {
    return res.status(400).send("POST Requests only");
  }
}
