import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from './ui/dialog';
import { Label } from './ui/label';
import { Input } from './ui/input';
import { useState } from 'react';

import { isValidEmail } from '@/lib/utils';
import { axiosInstance } from '@/lib/axiosInstance';
import { useAuth } from '@/lib/authCTX';
import { Button } from './ui/button';

export const AddFriendDialog = () => {
  const [open, setOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const [email, setEmail] = useState('');
  const [errorObj, setErrorObj] = useState({
    error: '',
    message: '',
    email: '',
  });

  const { user } = useAuth();

  const handleSubmit = async () => {
    if (!isValidEmail(email)) {
      setErrorObj((p) => ({ ...p, email: 'Email must be valid' }));
      return;
    }

    if (user?.email == email) {
      setErrorObj((p) => ({
        ...p,
        email: 'You can not be your own friend, buddy!',
      }));
      return;
    }

    setIsLoading(true);
    if (errorObj.email || errorObj.message || errorObj.error)
      setErrorObj({ email: '', message: '', error: '' });
    try {
      await axiosInstance.post('/friends', {
        email: email,
      });
      setErrorObj({ message: 'Added successfully!', error: '', email: '' });
      setTimeout(() => {
        setEmail('');
        setOpen(false);
        setErrorObj({ message: '', error: '', email: '' });
      }, 2000);
    } catch (err: any) {
      const { email, error, message } = err.response.data;
      setErrorObj({ email, error, message });
    } finally {
      setIsLoading(!true);
    }
  };

  return (
    <>
      <Dialog open={open}>
        <DialogTrigger asChild>
          <Button
            size='sm'
            onClick={() => {
              setOpen((p) => !p);
            }}
          >
            Add Friend +
          </Button>
        </DialogTrigger>
        <DialogContent className='sm:max-w-md'>
          <DialogHeader>
            <DialogTitle>Add Friend</DialogTitle>
          </DialogHeader>
          <div className='flex items-center '>
            <div className='grid flex-1 gap-2'>
              <Label htmlFor='link' className='sr-only'>
                Link
              </Label>
              <Input
                id='link'
                value={email}
                onChange={(e) => {
                  if (errorObj.email || errorObj.message || errorObj.error)
                    setErrorObj({ email: '', message: '', error: '' });

                  setEmail(e.target.value);
                }}
                placeholder='Enter email of the friend'
              />
            </div>
          </div>
          {!errorObj.error ? null : (
            <p className='text-xs italic text-muted text-red-400'>
              {errorObj.error}
            </p>
          )}
          {!errorObj.message ? null : (
            <p className='text-xs italic text-muted text-green-400'>
              {errorObj.message}
            </p>
          )}
          {!errorObj.email ? null : (
            <p className='text-xs italic text-muted text-red-400'>
              {errorObj.email}
            </p>
          )}

          <DialogFooter className='flex'>
            <DialogClose asChild>
              <Button
                type='button'
                variant='secondary'
                onClick={() => {
                  setOpen((p) => !p);
                }}
                disabled={isLoading}
              >
                Close
              </Button>
            </DialogClose>
            <Button
              type='button'
              variant='default'
              onClick={handleSubmit}
              disabled={isLoading}
            >
              Add
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
};
