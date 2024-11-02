import React, { useState } from "react";
import {
  Sidebar,
  Menu,
  MenuItem,
  SubMenu
} from "react-pro-sidebar";
import { Box, IconButton, Typography, useTheme } from "@mui/material";
import { Link } from "react-router-dom";
import { useLocation } from "react-router-dom";
import { UserRoleEnum, UserRoleName } from '../utils/userRoleEnum';

import HomeOutlinedIcon from "@mui/icons-material/HomeOutlined";
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import PaymentIcon from '@mui/icons-material/Payment';
import CorporateFareIcon from '@mui/icons-material/CorporateFare';
import BedroomChildIcon from '@mui/icons-material/BedroomChild';
import AccountBoxIcon from '@mui/icons-material/AccountBox';
import MenuOutlinedIcon from "@mui/icons-material/MenuOutlined";
import HistoryIcon from '@mui/icons-material/History';

const SideBar = () => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [toggled, setToggled] = useState(false);
  const [broken, setBroken] = useState(false);

  const theme = useTheme();
  const location = useLocation();

  const isActive = (path) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path + '/') || location.pathname === path;
  };

  const menuItems = [
    { path: "/dashboard", icon: <HomeOutlinedIcon />, text: "Dashboard", roles: [UserRoleEnum.ADMIN, UserRoleEnum.STAFF] },
    { path: "/booking", icon: <MenuBookIcon />, text: "Booking", roles: [UserRoleEnum.ADMIN, UserRoleEnum.STAFF] },
    {
      group: "Guest",
      icon: <PeopleAltIcon />,
      items: [
        { path: "/guest", text: "Guest", roles: [UserRoleEnum.ADMIN, UserRoleEnum.STAFF] },
        { path: "/guest-type", text: "Guest Types", roles: [UserRoleEnum.ADMIN, UserRoleEnum.STAFF] },
      ],
    },
    { path: "/payment", icon: <PaymentIcon />, text: "Payment", roles: [UserRoleEnum.ADMIN] },
    { path: "/rate_price", icon: <CorporateFareIcon />, text: "Rate Price", roles: [UserRoleEnum.ADMIN] },
    {
      group: "Room",
      icon: <BedroomChildIcon />,
      items: [
        { path: "/room", text: "Room", roles: [UserRoleEnum.ADMIN, UserRoleEnum.STAFF] },
        { path: "/room-type", text: "Room Types", roles: [UserRoleEnum.ADMIN, UserRoleEnum.STAFF] },
      ],
    },
    { path: "/user", icon: <AccountBoxIcon />, text: "User", roles: [UserRoleEnum.ADMIN] },
    { path: "/logs", icon: <HistoryIcon />, text: "Logs", roles: [UserRoleEnum.ADMIN] },
  ];

  const userRole = localStorage.getItem('userRole');

  return (
    <div style={{ display: "flex", height: "100%" }}>
      <Sidebar
        collapsed={isCollapsed}
        toggled={toggled}
        onBackdropClick={() => setToggled(false)}
        onBreakPoint={setBroken}
        breakPoint="md"
        backgroundColor="#f8f9fa"
        style={{
          height: "100%",
          borderRight: "1px solid #e0e0e0",
        }}
      >
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            height: "100%",
          }}
        >
          <Menu>
            {/* LOGO & MENU TOGGLE */}
            <MenuItem
              onClick={() => setIsCollapsed(!isCollapsed)}
              icon={isCollapsed ? <MenuOutlinedIcon /> : undefined}
              style={{
                margin: "10px 0 20px 0",
                color: "#333",
              }}
            >
              {!isCollapsed && (
                <Box
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                  ml="15px"
                >
                  <Typography variant="h6" color="#333" fontWeight="bold">
                    Hotel Management
                  </Typography>
                  <IconButton onClick={() => setIsCollapsed(!isCollapsed)}>
                    <MenuOutlinedIcon />
                  </IconButton>
                </Box>
              )}
            </MenuItem>

            {/* USER INFORMATION */}
            {!isCollapsed && (
              <Box mb="25px" textAlign="center">
                <Typography sx={{ color: "#333", fontWeight: "bold", mb: 1 }}>
                  {localStorage.getItem('username')}
                </Typography>
                <Typography sx={{ color: "#666", fontSize: "0.9rem" }}>
                  {UserRoleName[localStorage.getItem('userRole')]}
                </Typography>
              </Box>
            )}

            {/* MENU ITEMS */}
            {menuItems.map((item, index) => (
              item.group ? (
                <SubMenu
                  key={index}
                  icon={item.icon}
                  label={item.group}
                >
                  {item.items.map((subItem, subIndex) => (
                    (!subItem.roles || subItem.roles.includes(userRole)) && (
                      <MenuItem
                        key={`${index}-${subIndex}`}
                        component={<Link to={subItem.path} />}
                        active={isActive(subItem.path)}
                        style={{ backgroundColor: isActive(subItem.path) ? '#e0e0e0' : 'inherit' }}
                      >
                        {subItem.text}
                      </MenuItem>
                    )
                  ))}
                </SubMenu>
              ) : (
                (!item.roles || item.roles.includes(userRole)) && (
                  <MenuItem
                    key={index}
                    icon={item.icon}
                    component={<Link to={item.path} />}
                    active={isActive(item.path)}
                    style={{ backgroundColor: isActive(item.path) ? '#e0e0e0' : 'inherit' }}
                  >
                    {item.text}
                  </MenuItem>
                )
              )
            ))}
          </Menu>

          {/* SIDEBAR FOOTER */}
          <Box mt="auto" textAlign="center" py={2}>
            <Typography variant="caption" color="#666">
              © 2024 Hotel Management
            </Typography>
          </Box>
        </Box>
      </Sidebar>

      {/* MAIN CONTENT AREA */}
      <main style={{ flex: 1, padding: "16px" }}>
        {broken && (
          <IconButton onClick={() => setToggled(!toggled)}>
            <MenuOutlinedIcon />
          </IconButton>
        )}
        {/* Content goes here */}
      </main>
    </div>
  );
};

export default SideBar;
