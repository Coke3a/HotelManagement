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
      justifyContent="flex-end"
      alignItems="center"
      p={2}
      bgcolor="var(--background-color)"
      boxShadow="0 2px 4px rgba(0, 0, 0, 0.1)"
      borderBottom="1px solid var(--secondary-color)"
    >
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
