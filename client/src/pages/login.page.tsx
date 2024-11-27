import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

import { useAuth } from '@/lib/authCTX';
import { axiosInstance } from '@/lib/axiosInstance';
import { useState } from 'react';
import { Link } from 'react-router-dom';

export function LoginPage() {
  const { setUser } = useAuth();

  const [userData, setUserData] = useState({
    email: `user1@gmail.com`,
    password: '123456',
  });

  const [userError, setUserError] = useState({
    email: '',
    password: '',
  });
  const handleOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const tgt = e.target;

    // @ts-ignore
    if (userError[tgt.name]) {
      setUserError((p) => ({ ...p, [tgt.name]: '' }));
    }

    setUserData((p) => ({ ...p, [tgt.name]: tgt.value }));
  };
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async () => {
    try {
      setIsLoading(true);
      if (!userData.email) {
        setUserError((p) => ({ ...p, email: 'Email is required' }));
        return;
      }
      if (!userData.password) {
        setUserError((p) => ({ ...p, password: 'password is required' }));
        return;
      }

      const { data } = await axiosInstance.post('/auth/login', userData);

      setUser(data.user, data.jwt);
    } catch (error: any) {
      console.error(
        'error while hitting signup endpoint: reason',
        error.response.data
      );

      const { email, password } = error.response.data;
      setUserError({ email, password });
    } finally {
      setIsLoading(false);
    }
  };
  return (
    <Card className='w-[350px] h-fit'>
      <CardHeader>
        <CardTitle className='text-3xl text-center'>Login</CardTitle>
      </CardHeader>
      <CardContent>
        <form>
          <div className='grid w-full items-center gap-4'>
            <div className='flex flex-col space-y-1.5'>
              <Label htmlFor='email'>Email</Label>
              <Input
                onChange={handleOnChange}
                value={userData.email}
                type='email'
                id='email'
                name='email'
                placeholder='name@email.com'
              />
              <p className='text-xs italic text-muted text-red-400'>
                {userError.email}
              </p>
            </div>
            <div className='flex flex-col space-y-1.5'>
              <Label htmlFor='password'>password</Label>
              <Input
                onChange={handleOnChange}
                value={userData.password}
                name='password'
                type='password'
                id='password'
                placeholder='12345678'
              />
              <p className='text-xs italic text-muted text-red-400'>
                {userError.password}
              </p>
            </div>
          </div>
        </form>
      </CardContent>
      <CardFooter className='flex flex-col'>
        <div className='text-xs mb-4'>
          <p>
            Does not have an account?
            <Link to='/auth/signup' className='text-blue-600 underline'>
              Signup
            </Link>
          </p>
        </div>
        <Button size='sm' onClick={handleSubmit} disabled={isLoading}>
          Submit
        </Button>
      </CardFooter>
    </Card>
  );
}
