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

import { Link } from 'react-router-dom';
import axios from 'axios';
import { useState } from 'react';
import { useAuth } from '@/lib/authCTX';
import { axiosInstance } from '@/lib/axiosInstance';

export function SignupPage() {
  const { setUser } = useAuth();

  const [userData, setUserData] = useState({
    email: `imsamad00@gmail.com`,
    password: 'pwd@Hello123',
    name: 'Abdus Samad',
  });

  const [userError, setUserError] = useState({
    email: '',
    password: '',
    name: '',
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

      console.log(userData);

      const { data } = await axiosInstance.post('/auth/signup', userData);
      setUser(data.user);
    } catch (error: any) {
      console.error(
        'error while hitting signup endpoint: reason',
        error.response.data
      );

      const { email, name, password } = error.response.data;
      setUserError({ email, name, password });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className='w-[350px] h-fit'>
      <CardHeader>
        <CardTitle className='text-3xl text-center'>Signup</CardTitle>
      </CardHeader>
      <CardContent>
        <form>
          <div className='grid w-full items-center gap-4'>
            <div className='flex flex-col space-y-1.5'>
              <Label htmlFor='name'>Email</Label>
              <Input
                onChange={handleOnChange}
                value={userData.name}
                type='text'
                name='name'
                id='name'
                placeholder='name@email.com'
              />
              <p className='text-xs italic text-muted text-red-400'>
                {userError.name}
              </p>
            </div>
            <div className='flex flex-col space-y-1.5'>
              <Label htmlFor='email'>Email</Label>
              <Input
                onChange={handleOnChange}
                value={userData.email}
                type='email'
                name='email'
                id='email'
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
                type='password'
                name='password'
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
            Already have an account?
            <Link to='/auth/login' className='text-blue-600 underline'>
              Login
            </Link>
          </p>
        </div>
        <Button onClick={handleSubmit} size='sm' disabled={isLoading}>
          Submit
        </Button>
      </CardFooter>
    </Card>
  );
}
