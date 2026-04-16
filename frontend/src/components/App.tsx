import { Outlet } from "react-router";

export default function App() {
  return (
    <>
    <div className="flex w-screen h-screen justify-center items-center">
      <Outlet />
    </div>
    </>
  )
}
