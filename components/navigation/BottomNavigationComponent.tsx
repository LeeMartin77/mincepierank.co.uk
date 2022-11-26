import { BottomNavigation, BottomNavigationAction } from "@mui/material";
import Link from "next/link";

import { navigationConfig } from "./NavigationConfig";

export function BottomNavigationComponent() {
  return (
    <BottomNavigation
      showLabels
      style={{
        position: "fixed",
        bottom: 0,
        width: "100%",
      }}
    >
      {navigationConfig.map((x, i) => (
        <BottomNavigationAction
          LinkComponent={Link}
          href={x.route}
          key={i}
          icon={<x.icon />}
          label={x.label}
        />
      ))}
    </BottomNavigation>
  );
}
