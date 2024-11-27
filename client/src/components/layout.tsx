import { Outlet } from 'react-router-dom';
import { Navbar } from './navbar/navbar';

export const Layout = () => {
  return (
    <div className='bg-gray-900'>
      <div className='h-screen min-h-screen max-h-screen max-w-screen-xl w-full overflow-hidden flex flex-col container mx-auto '>
        <Navbar />
        <div className='flex-1 flex justify-center green'>
          <Outlet />
        </div>
      </div>
    </div>
  );
};
