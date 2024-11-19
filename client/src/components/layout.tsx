import { Outlet } from 'react-router-dom';
import { Navbar } from './navbar/navbar';

export const Layout = () => {
  return (
    <div className='h-screen min-h-screen flex flex-col '>
      <Navbar />
      <div className='flex-1 flex justify-center p-10'>
        <Outlet />
      </div>
    </div>
  );
};
