import { Route, Routes } from "react-router";
import HomePage from "./default-components/HomePage";
import Login from "./default-components/Login";

export default function App() {
  return (
    <>
    <div className="flex w-screen h-screen justify-center items-center">
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </div>
    </>
  )
}
