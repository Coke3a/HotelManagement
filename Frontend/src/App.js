import React from "react";
import PropTypes from "prop-types";
import HeaderBar from "./layout/HeaderBar";
import SideBar from "./layout/SideBar";

import { Routes, Route, useLocation, Navigate } from "react-router-dom";
import { CssBaseline, Box } from "@mui/material";

import Booking from "./components/Booking";
import BookingAdd from "./components/BookingAdd";
import BookingEdit from "./components/BookingEdit";
import DashBoard from "./components/DashBoard";
import Guest from "./components/Guest";
import GuestAdd from "./components/GuestAdd";
import GuestEdit from "./components/GuestEdit";
import Payment from "./components/Payment";
import PaymentAdd from "./components/PaymentAdd";
import PaymentEdit from "./components/PaymentEdit";
import RatePrice from "./components/RatePrice";
import RatePriceAdd from "./components/RatePriceAdd";
import RatePriceEdit from "./components/RatePriceEdit";
import Room from "./components/Room";
import RoomAdd from "./components/RoomAdd";
import RoomEdit from "./components/RoomEdit";
import User from "./components/User";
import UserAdd from "./components/UserAdd";
import UserEdit from "./components/UserEdit";
import Login from "./components/auth/Login";
import ProtectedRoute from "./components/auth/ProtectedRoute";
import RoomType from "./components/RoomType";
import RoomTypeAdd from "./components/RoomTypeAdd";
import RoomTypeEdit from "./components/RoomTypeEdit";
import GuestType from "./components/GuestType";
import GuestTypeAdd from "./components/GuestTypeAdd";
import GuestTypeEdit from "./components/GuestTypeEdit";
import ReceiptPrint from './components/ReceiptPrint';
import Log from "./components/Log";
import DailyBookingSummary from './components/DailyBookingSummary';
import DailyBookingSummaryEdit from './components/DailyBookingSummaryEdit';
import { logout } from './utils/auth';
import { UserRoleEnum } from './utils/userRoleEnum';

export { logout };

const withAuth = (Component) => {
  return (props) => {
    const token = localStorage.getItem('token');
    if (!token) {
      return <Navigate to="/login" replace />;
    }
    return <Component {...props} />;
  };
};

const App = () => {
  const location = useLocation();
  const isLoginPage = location.pathname === "/login";

  return (
    <>
      <CssBaseline />
      <div className="app">
        {!isLoginPage && <SideBar />}
        <main className="content">
          {!isLoginPage && <HeaderBar />}
          <div className="content_body">
            <Box m="20px">
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/unauthorized" element={<div>Unauthorized Access</div>} />
                <Route element={<ProtectedRoute allowedRoles={[UserRoleEnum.ADMIN, UserRoleEnum.STAFF]} />}>
                  <Route path="/" element={<Navigate to="/dashboard" replace />} />
                  <Route path="/dashboard" element={<DashBoard />} />
                  <Route path="/booking" element={<Booking />} />
                  <Route path="/booking/add" element={<BookingAdd />} />
                  <Route path="/booking/edit/:id" element={<BookingEdit />} />
                  <Route path="/guest" element={<Guest />} />
                  <Route path="/guest/add" element={<GuestAdd />} />
                  <Route path="/guest/edit/:id" element={<GuestEdit />} />
                  <Route path="/guest-type" element={<GuestType />} />
                  <Route path="/room" element={<Room />} />
                  <Route path="/room-type" element={<RoomType />} />
                </Route>
                <Route element={<ProtectedRoute allowedRoles={[UserRoleEnum.ADMIN]} />}>
                  <Route path="/payment" element={<Payment />} />
                  <Route path="/payment/add" element={<PaymentAdd />} />
                  <Route path="/payment/edit/:id" element={<PaymentEdit />} />
                  <Route path="/rate_price" element={<RatePrice />} />
                  <Route path="/rate_price/add" element={<RatePriceAdd />} />
                  <Route path="/rate_price/edit/:id" element={<RatePriceEdit />} />
                  <Route path="/room/add" element={<RoomAdd />} />
                  <Route path="/room/edit/:id" element={<RoomEdit />} />
                  <Route path="/user" element={<User />} />
                  <Route path="/user/add" element={<UserAdd />} />
                  <Route path="/user/edit/:id" element={<UserEdit />} />
                  <Route path="/room-type/add" element={<RoomTypeAdd />} />
                  <Route path="/room-type/edit/:id" element={<RoomTypeEdit />} />
                  <Route path="/guest-type/add" element={<GuestTypeAdd />} />
                  <Route path="/guest-type/edit/:id" element={<GuestTypeEdit />} />
                  <Route path="/receipt/:id" element={<ReceiptPrint />} />
                  <Route path="/logs" element={<Log />} />
                  <Route path="/daily-summary" element={<DailyBookingSummary />} />
                  <Route path="/daily-summary/edit/:date" element={<DailyBookingSummaryEdit />} />
                </Route>
              </Routes>
            </Box>
          </div>
        </main>
      </div>
    </>

  )
}

export default App;
