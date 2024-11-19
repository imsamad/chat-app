import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  RouterProvider,
} from 'react-router-dom';
import { HomaPage } from './pages/home.page';
import { SignupPage } from './pages/signup.page';
import { LoginPage } from './pages/login.page';
import { Layout } from './components/layout';

const routers = createBrowserRouter(
  createRoutesFromElements(
    <Route path='/' element={<Layout />}>
      <Route path='/' element={<HomaPage />} />
      <Route path='/auth/login' element={<LoginPage />} />
      <Route path='/auth/signup' element={<SignupPage />} />
    </Route>
  )
  //   [
  //   {
  //     path: '/',
  //     element: <HomaPage />,
  //   },
  //   {
  //     path: '/auth/login',
  //     element: <LoginPage />,
  //   },
  //   {
  //     path: '/auth/signup',
  //     element: <SignupPage />,
  //   },
  // ]
);

export const AppRoutes = () => <RouterProvider router={routers} />;
