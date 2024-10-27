export const handleTokenExpiration = (error, navigate) => {
  if (error.message === "access token has expired") {
    localStorage.clear();
    navigate('/login', { 
      state: { 
        error: "Your session has expired. Please login again." 
      } 
    });
  }
};