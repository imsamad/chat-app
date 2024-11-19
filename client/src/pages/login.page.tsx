import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Link } from 'react-router-dom';

export function LoginPage() {
  return (
    <Card className='w-[350px] h-fit'>
      <CardHeader>
        <CardTitle className='text-3xl text-center'>Login</CardTitle>
      </CardHeader>
      <CardContent>
        <form>
          <div className='grid w-full items-center gap-4'>
            <div className='flex flex-col space-y-1.5'>
              <Label htmlFor='name'>Email</Label>
              <Input type='email' id='email' placeholder='name@email.com' />
            </div>
            <div className='flex flex-col space-y-1.5'>
              <Label htmlFor='password'>password</Label>
              <Input type='password' id='password' placeholder='12345678' />
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
        <Button size='sm'>Submit</Button>
      </CardFooter>
    </Card>
  );
}
