import Cookies from 'js-cookie';

// App
const sidebarStatusKey = 'sidebar_status';
export const getSidebarStatus = () => Cookies.get(sidebarStatusKey);
export const setSidebarStatus = (sidebarStatus: string) => Cookies.set(sidebarStatusKey, sidebarStatus);

// User
const cookieKey = 'pfcookie';
export const cookiePresent = (): boolean => {
  const cookie = Cookies.get(cookieKey);
  return (cookie !== undefined && cookie !== null && cookie !== '');
};
