export const getToken = () => localStorage.getItem('token');

export const setToken = (token) => localStorage.setItem('token', token);

export const removeToken = () => localStorage.removeItem('token');

export const getUserRole = () => localStorage.getItem('userRole');

export const setUserRole = (role) => localStorage.setItem('userRole', role);

export const getUserId = () => localStorage.getItem('userId');

export const setUserId = (id) => localStorage.setItem('userId', id);

export const isAuthenticated = () => !!getToken();

export const logout = () => {
  removeToken();
  localStorage.removeItem('userRole');
  localStorage.removeItem('userId');
};
export const getUsername = () => localStorage.getItem('username');

export const setUsername = (username) => localStorage.setItem('username', username);

// export const handleTokenExpiration = (error, navigate) => {
//   if (error.message.includes("access token has expired")) {
//     // Clear all auth data
//     localStorage.removeItem('token');
//     localStorage.removeItem('userRole');
//     localStorage.removeItem('username');
    
//     // Redirect to login with error message
//     navigate('/login', {
//       state: { error: 'Your session has expired. Please log in again.' }
//     });
//   }
// };