import {
  createBrowserRouter,
  createRoutesFromElements,
  Outlet,
  redirect,
  Route,
  RouterProvider,
  useNavigate,
} from 'react-router-dom';
import { HomaPage } from './pages/home.page';
import { SignupPage } from './pages/signup.page';
import { LoginPage } from './pages/login.page';
import { Layout } from './components/layout';
import { useAuth } from './lib/authCTX';
import { useEffect } from 'react';
import { ChatPage } from './pages/chat.page';

const GuestPages = () => {
  const { isLoggedIn } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isLoggedIn === true) navigate('/chat');
  }, [isLoggedIn]);
  return <Outlet />;
};

const ProtectedPages = () => {
  const { isLoggedIn } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isLoggedIn === undefined) return;
    if (isLoggedIn === false) navigate('/auth/login');
    if (isLoggedIn == true) navigate('/chat');
  }, [isLoggedIn]);
  return <Outlet />;
};

const routers = createBrowserRouter(
  createRoutesFromElements(
    <Route path='/' element={<Layout />}>
      <Route element={<GuestPages />}>
        <Route index element={<HomaPage />} />
        <Route path='auth/login' element={<LoginPage />} />
        <Route path='auth/signup' element={<SignupPage />} />
      </Route>
      <Route element={<ProtectedPages />}>
        <Route path='/chat' element={<ChatPage />} />
      </Route>
    </Route>
  )
);

export const AppRoutes = () => <RouterProvider router={routers} />;
