// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from "next";
import { getToken } from "next-auth/jwt";
import {
  addPieRanking,
  getPieByMakerAndId,
  getMyRankingForPie,
  StorageError,
} from "../../../../../system/storage";
import { timeInSeason } from "../../../../../system/seasonConfig";

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse,
  fnGetPie = getPieByMakerAndId,
  fnAddPieRanking = addPieRanking,
  fnGetMyRankingForPie = getMyRankingForPie,
  fnGetToken = getToken,
  fnGetTime = Date.now
) {
  if (!["GET", "POST"].includes(req.method as string)) {
    return res.status(400).send("GET or POST Requests only");
  }
  return fnGetToken({ req }).then((token) => {
    if (!token) {
      res.status(401).send("Requires Authentication");
      return;
    }
    const userid = (token.id ?? token.email ?? token.name) as string;
    if (req.method === "GET") {
      const { brandid, pieid } = req.query;
      return fnGetMyRankingForPie(
        brandid as string,
        pieid as string,
        userid
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
      if (!timeInSeason(fnGetTime())) {
        return res.status(400).send("Out of Season");
      }
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
            fnAddPieRanking({ ...parsedBody, userid }).then((sres) => {
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
  });
}
