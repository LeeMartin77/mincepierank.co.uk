// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from "next";
import {
  addPieRanking,
  getPieByMakerAndId,
  getMyRankingForPie,
  StorageError,
} from "../../../../../system/storage";

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse,
  fnGetPie = getPieByMakerAndId,
  fnAddPieRanking = addPieRanking,
  fnGetMyRankingForPie = getMyRankingForPie
) {
  if (req.method === "GET") {
    const { brandid, pieid } = req.query;
    return fnGetMyRankingForPie(
      brandid as string,
      pieid as string,
      req.headers["userid"] as string
    ).then((response) => {
      response
        .map((ranking) => {
          res.status(200).send(ranking);
        })
        .mapErr((se) => {
          if (se === StorageError.NotFound) {
            res.status(404).send("Not Found");
          } else {
            res.status(500).send("Server Error");
          }
        });
    });
  }
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
    return res.status(400).send("GET or POST Requests only");
  }
}
