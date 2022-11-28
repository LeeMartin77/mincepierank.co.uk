import {
  Button,
  Card,
  CardActions,
  CardMedia,
  Grid,
  Link,
} from "@mui/material";
import { Maker } from "../system/storage";
import styles from "./brandCard.module.css";

export function BrandCard({ maker }: { maker: Maker }) {
  return (
    <Grid key={maker.id} item xs={6} sm={4} md={3}>
      <Card>
        <Link href={`/brands/${maker.id}`}>
          <CardMedia
            className={styles.brandImage}
            component="img"
            height="150"
            image={maker.logo}
            alt={`${maker.name} Logo`}
          />
        </Link>
        <CardActions>
          <Button
            LinkComponent={Link}
            href={`/brands/${maker.id}`}
            className={styles.brandButton}
          >
            {maker.name}
          </Button>
        </CardActions>
      </Card>
    </Grid>
  );
}
