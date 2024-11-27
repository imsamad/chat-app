import { Link } from 'react-router-dom';
import { Button } from '../ui/button';
import { useAuth } from '@/lib/authCTX';
import { AddFriendDialog } from '../AddFriendDialog';

export function Navbar() {
  const { user, logout } = useAuth();

  const isLoggedIn = !!user;

  return (
    <div className='flex p-4  mx-auto items-center border-b border-gray-700 w-full log'>
      <Link to='/'>
        <h1 className='text-4xl font-bold text-gray-100'>Logo</h1>
      </Link>
      <div className='flex flex-1 items-center justify-end gap-2'>
        {isLoggedIn ? (
          <>
            <AddFriendDialog />
            <Button size='sm' color='danger' onClick={logout}>
              Logout
            </Button>
          </>
        ) : (
          <>
            <Link to='/auth/login'>
              <Button>Login</Button>
            </Link>
            <Link to='/auth/signup'>
              <Button>Signup</Button>
            </Link>
          </>
        )}
      </div>
    </div>
  );
}
