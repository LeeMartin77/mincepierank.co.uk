import { MakerPie, PieRankingSummary } from "../../system/storage";
import { Button, Card, CardActions, CardMedia, Link } from "@mui/material";

export function PieSummaryLink({ pie }: { pie: MakerPie }) {
  return (
    <Card data-testid={pie.id + "-link-card"}>
      <Link href={`/brands/${pie.makerid}/${pie.id}`}>
        <CardMedia
          component="img"
          height="200"
          image={pie.image_file}
          alt={`${pie.displayname}`}
        />
      </Link>
      <CardActions>
        <Button
          LinkComponent={Link}
          href={`/brands/${pie.makerid}/${pie.id}`}
          style={{ width: "100%", textAlign: "center" }}
        >
          {pie.displayname}
        </Button>
      </CardActions>
    </Card>
  );
}

export function PieList({
  pies,
  rankings,
}: {
  pies: MakerPie[];
  rankings: PieRankingSummary[];
}) {
  return (
    <>
      {pies.map((pie) => (
        <PieSummaryLink key={pie.id} pie={pie} />
      ))}
    </>
  );
}
