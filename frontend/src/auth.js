import axiosInstance from './utils/axiosInstance';

export const signup = (email, password) => {
  return axiosInstance.post('/api/signup', { email, password });
};

export const login = (email, password) => {
  return axiosInstance.post('/api/login', { email, password });
};