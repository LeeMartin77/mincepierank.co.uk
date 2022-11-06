import { Button, Card, CardActions, CardMedia, Link } from "@mui/material";
import { MakerPie } from "../system/storage";

export function PieSummaryLink({ pie }: { pie: MakerPie }) {
  return (
    <Card>
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
