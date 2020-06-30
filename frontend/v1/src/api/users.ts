import request from '@/utils/request';
import { LoginRequest } from '@/models/state';
import { methods } from '@/constants/http';

export const getUserInfo = () =>
  request({
    url: '/users/status',
    method: methods.GET,
    withCredentials: true,
  });

export const login = (data: LoginRequest) =>
  request({
    url: '/users/login',
    method: methods.POST,
    data: data,
  });

export const logout = () =>
  request({
    url: '/users/logout',
    method: methods.POST,
  });
