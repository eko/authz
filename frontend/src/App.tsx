import {
  createBrowserRouter,
  RouterProvider,
} from 'react-router-dom';
import { ToastProvider } from 'context/toast';
import { routes } from 'routes';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterMoment } from '@mui/x-date-pickers/AdapterMoment';

import './App.css'

export default function App() {
  const router = createBrowserRouter(routes);

  return (
    <LocalizationProvider dateAdapter={AdapterMoment}>
      <ToastProvider>
        <RouterProvider router={router} />
      </ToastProvider>
    </LocalizationProvider>
  );
}