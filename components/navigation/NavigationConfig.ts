import { FormatListBulleted, Home, Person } from "@mui/icons-material";

export const navigationConfig = [
  {
    label: "Home",
    route: "/",
    icon: Home,
  },
  {
    label: "All Pies",
    route: "/all-mince-pies",
    icon: FormatListBulleted,
  },
  {
    label: "My Rankings",
    route: "/my-rankings",
    icon: Person,
  },
];
