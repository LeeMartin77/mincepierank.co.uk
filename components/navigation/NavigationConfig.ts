import { FormatListBulleted, Shop, Info, Home } from "@mui/icons-material";

export const navigationConfig = [
  {
    label: "Home",
    route: "/",
    icon: Home,
    mobileHidden: false,
  },
  {
    label: "List",
    route: "/all-mince-pies",
    icon: FormatListBulleted,
    mobileHidden: false,
  },
  {
    label: "Brands",
    route: "/brands",
    icon: Shop,
    mobileHidden: true,
  },
  {
    label: "About",
    route: "/about",
    icon: Info,
    mobileHidden: true,
  },
];
