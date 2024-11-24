import React, { useState } from "react";
import { Box, IconButton, Menu, MenuItem, Typography } from "@mui/material";
import PersonOutlinedIcon from "@mui/icons-material/PersonOutlined";
import { Link, useNavigate } from "react-router-dom";
import { logout } from "../utils/auth";

const HeaderBar = () => {
  const [anchorEl, setAnchorEl] = useState(null);
  const navigate = useNavigate();

  const handleMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <Box
      display="flex"
      justifyContent="space-between"
      alignItems="center"
      p={1}
      bgcolor="var(--background-color)"
      boxShadow="0 2px 4px rgba(0, 0, 0, 0.1)"
      borderBottom="1px solid var(--secondary-color)"
    >
      {/* Page Title Section */}
      <Typography
        variant="h5"
        sx={{
          fontWeight: 500,
          color: '#1a1a1a'
        }}
      >
        {/* Dynamic title based on current route */}
        {window.location.pathname === '/user' && 'User Management'}
        {window.location.pathname === '/room' && 'Room Management'}
        {window.location.pathname === '/room-type' && 'Room Type'}
        {window.location.pathname === '/rate_price' && 'Rate Price'}
        {window.location.pathname === '/guest' && 'Guest'}
        {window.location.pathname === '/guest-type' && 'Guest Type'}
        {window.location.pathname === '/logs' && 'System Logs'}
        {window.location.pathname === '/daily-summary' && 'Daily Booking Summary'}
      </Typography>

      {/* Icons Section */}
      <Box display="flex" alignItems="center">
        <IconButton
          sx={{
            color: "#4299e1",
            border: "1px solid #4299e1",
            borderRadius: "50%",
            padding: "8px",
            backgroundColor: "white",
            '&:hover': {
              backgroundColor: "#e0f7fa",
            },
          }}
          onClick={handleMenu}
        >
          <PersonOutlinedIcon />
        </IconButton>

        {/* Dropdown Menu */}
        <Menu
          id="menu-appbar"
          anchorEl={anchorEl}
          anchorOrigin={{
            vertical: "top",
            horizontal: "right",
          }}
          keepMounted
          transformOrigin={{
            vertical: "top",
            horizontal: "right",
          }}
          open={Boolean(anchorEl)}
          onClose={handleClose}
        >
          <Link to="/profile" className="menu-bars">
            <MenuItem onClick={handleClose}>
              <Typography color="#1a365d">Profile</Typography>
            </MenuItem>
          </Link>
          <MenuItem onClick={handleLogout}>
            <Typography color="#1a365d">Logout</Typography>
          </MenuItem>
        </Menu>
      </Box>
    </Box>
  );
};

export default HeaderBar;
