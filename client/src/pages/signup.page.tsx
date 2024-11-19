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
import { FormEventHandler, useState } from 'react';

export function SignupPage() {
  const [userData, setUserData] = useState({
    email: 'imsamad00@gmail.com',
    password: 'pwd@Hello123',
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

      console.log(userData);

      const { data } = await axios.post(
        'http://localhost:4000/auth/signup',
        userData
      );

      console.log(data);
    } catch (error) {
      console.error('error while hitting signup endpoint: reason', error);
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
                value={userData.email}
                type='email'
                name='email'
                id='email'
                placeholder='name@email.com'
              />
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
            </div>
          </div>
        </form>
      </CardContent>
      <CardFooter className='flex flex-col'>
        <div className='text-xs mb-4'>
          <p>
            Already have an account?
            <Link to='/auth/login' className='text-blue-600 underline'>
              Signup
            </Link>
          </p>
        </div>
        <Button onClick={handleSubmit} size='sm'>
          Submit
        </Button>
      </CardFooter>
    </Card>
  );
}
