import {
  BottomNavigation,
  BottomNavigationAction,
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
} from "@mui/material";
import Link from "next/link";
import MenuIcon from "@mui/icons-material/Menu";

import { navigationConfig } from "./NavigationConfig";
import { useState } from "react";

export function BottomNavigationComponent() {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <>
      <BottomNavigation
        showLabels
        style={{
          position: "fixed",
          bottom: 0,
          width: "100%",
        }}
      >
        {navigationConfig
          .filter((x) => !x.mobileHidden)
          .map((x, i) => (
            <BottomNavigationAction
              LinkComponent={Link}
              href={x.route}
              key={i}
              icon={<x.icon />}
              label={x.label}
            />
          ))}
        <BottomNavigationAction
          icon={<MenuIcon />}
          onClick={handleClick}
          label="More"
        />
      </BottomNavigation>
      <Menu
        id="more-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        MenuListProps={{
          "aria-labelledby": "basic-button",
        }}
      >
        {navigationConfig
          .filter((x) => x.mobileHidden)
          .map((x, i) => {
            return (
              <MenuItem
                key={i}
                onClick={handleClose}
                LinkComponent={Link}
                href={x.route}
              >
                <ListItemIcon>
                  <x.icon />
                </ListItemIcon>
                <ListItemText>{x.label}</ListItemText>
              </MenuItem>
            );
          })}
      </Menu>
    </>
  );
}
