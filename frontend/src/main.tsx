import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './css/index.css'
import { createBrowserRouter, RouterProvider } from 'react-router'
import HomePage from './components/default-components/HomePage.tsx'
import Login from './components/default-components/Login.tsx'
import App from './components/App.tsx'
import StorageHome from './components/auth-components/storage-components/StorageHome.tsx'
import { authMiddleware, loginMiddleware } from './middleware.ts'

const router = createBrowserRouter([
  {
    path: "/", 
    element: <App />,
    children: [
      {index: true, element: <HomePage />},
      {
        middleware: [loginMiddleware],
        path: "login",
        element: <Login />
      },
      {
        path: "storage",
        element: <App />,
        middleware: [authMiddleware],
        children: [
          {index: true, element: <StorageHome />},
        ]
      },
    ]
  },
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <RouterProvider router={router} />
  </StrictMode>,
)
