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