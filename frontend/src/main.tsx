import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { createBrowserRouter, RouterProvider } from 'react-router'
import HomePage from './components/default-components/HomePage.tsx'
import Login from './components/default-components/Login.tsx'
import App from './components/App.tsx'
import StorageHome from './components/auth-components/storage-components/StorageHome.tsx'
import { authMiddleware } from './middleware.ts'

const router = createBrowserRouter([
  {
    path: "/", 
    element: <App />,
    children: [
      {index: true, element: <HomePage />},
      {path: "login", element: <Login />},
      {
        middleware: [authMiddleware],
        path: "storage", 
        element: <StorageHome />,
      },
    ]
  },
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
)
