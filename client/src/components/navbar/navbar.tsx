import { Button } from '../ui/button';

export function Navbar() {
  return (
    <div className='flex p-4'>
      <h1>Logo</h1>
      <div className='flex flex-1 justify-end'>
        <Button>Login</Button>
      </div>
    </div>
  );
}
